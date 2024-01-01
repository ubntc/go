package kafkago

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/segmentio/kafka-go"
)

var reKafkaErrorCode = regexp.MustCompile(`^\[[0-9]+\]`)

var (
	ErrOffsetRequired = errors.New("offset required")
	ErrInvalidOffset  = errors.New("invalid offset")
)

func KafkaError(err error) kafka.Error {
	if err == nil {
		return 0
	}

	match := reKafkaErrorCode.FindString(err.Error())
	if len(match) < 3 {
		return kafka.Unknown
	}
	code, _ := strconv.Atoi(match[1 : len(match)-1])
	return kafka.Error(code)
}
