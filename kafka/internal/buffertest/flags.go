package buffertest

import (
	"flag"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func LoadConfig() Config {
	var broker = flag.String("b", "localhost:9092", "Kafka brokers (comma-separated)")
	var group = flag.String("G", "buffertest.tester", "Kafka consumer group")
	var topic = flag.String("t", "buffertest", "Kafka topic")
	var timeout = flag.Duration("d", time.Second*20, "send and receive timeout")
	var numEvents = flag.Int("c", 10, "number of events to send and receive")
	var numPartitions = flag.Int("p", 10, "number of partitions used when creating topics")

	flag.Parse()
	brokers := strings.Split(*broker, ",")
	return Config{
		NumEvents:       *numEvents,
		PipelineTimeout: *timeout,
		WriterTick:      *timeout / time.Duration(*numEvents) / 10,
		WaiterTick:      *timeout / 10,
		Writer: kafka.WriterConfig{
			Brokers:      brokers,
			Topic:        *topic,
			BatchTimeout: time.Millisecond,
		},
		Reader: kafka.ReaderConfig{
			Brokers:       brokers,
			Topic:         *topic,
			GroupID:       *group,
			QueueCapacity: 1,
		},
		Topic: kafka.TopicConfig{
			Topic:             *topic,
			NumPartitions:     *numPartitions,
			ReplicationFactor: 1,
		},
	}
}
