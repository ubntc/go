package batbq

import (
	"context"

	"cloud.google.com/go/bigquery"
)

// confirmMessages acks and nacks `messages` in the context of a potential
// batching `error` and returns the number of acked and nacked messages.
func confirmMessages(messages []Message, err error) (numAcked int, numNacked int) {
	nacked := handleErrors(messages, err)

	switch {
	case len(nacked) == len(messages):
		// all messages had errors and are already nacked
	case len(nacked) == 0:
		// no messages had errors and can be acked
		for _, m := range messages {
			m.Ack()
		}
	default:
		// some messages had errors, we need to check which
		for i, m := range messages {
			if _, ok := nacked[i]; ok {
				continue
			}
			m.Ack()
		}
	}
	return len(messages) - len(nacked), len(nacked)
}

// handleErrors nacks `messages` according to the type of the received `error`.
// It returns an index of the nacked messages.
func handleErrors(messages []Message, err error) (index map[int]struct{}) {
	if err == nil {
		return nil
	}
	nacked := make(map[int]struct{})
	mulErr, isMulti := err.(bigquery.PutMultiError)
	switch {
	case isMulti:
		for _, insErr := range mulErr {
			messages[insErr.RowIndex].Nack(insErr.Errors)
			nacked[insErr.RowIndex] = struct{}{}
		}
	case err == context.Canceled:
		// batcher is shutdown down, just nack the messages without forwarding the error
		for i, m := range messages {
			m.Nack(nil)
			nacked[i] = struct{}{}
		}
	default:
		// another error happened, forward it with the Nack to allow handling upstream
		for i, m := range messages {
			m.Nack(err)
			nacked[i] = struct{}{}
		}
	}
	return nacked
}
