package buffertest

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Writer          kafka.WriterConfig
	Reader          kafka.ReaderConfig
	Topic           kafka.TopicConfig
	NumEvents       int
	PipelineTimeout time.Duration
	WriterTick      time.Duration
	WaiterTick      time.Duration
}
