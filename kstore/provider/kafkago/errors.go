package kafkago

import (
	"log"
	"regexp"
	"strconv"

	"github.com/segmentio/kafka-go"
	"github.com/ubntc/go/kstore/kstore"
)

func NewLogger(name string) kstore.LoggerFunc {
	return func(format string, args ...any) {
		log.Printf(name+": "+format, args...)
		log.Println()
	}
}

func NilLogger() kstore.LoggerFunc {
	return func(format string, args ...any) {}
}

var reKafkaErrorCode = regexp.MustCompile(`^\[[0-9]+\]`)

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
