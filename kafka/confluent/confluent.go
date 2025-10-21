package confluent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xanderxampp-be/franco/contextwrap"
	"github.com/xanderxampp-be/franco/trace"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.elastic.co/apm"
)

type KafkaConfig struct {
	Brokers         []string
	Topic           string
	ConsumerGroup   string
	KafkaUsr        string
	KafkaPsw        string
	KafkaMechanisms string
	KafkaProtocol   string
	KafkaLinger     string
	KafkaBatch      string
	KafkaFlush      string
	KafkaGroup      string
	KafkaOffset     string
	KafkaSession    string
	KafkaHeartbeat  string
	KafkaFetch      string
	KafkaMaxpoll    string
}

func (r KafkaConfig) NewConfluentReader() *kafka.Consumer {
	brokersJoined := strings.Join(r.Brokers, `,`)
	read, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     brokersJoined,
		"sasl.mechanisms":       r.KafkaMechanisms,
		"security.protocol":     r.KafkaProtocol,
		"sasl.username":         r.KafkaUsr,
		"sasl.password":         r.KafkaPsw,
		"group.id":              r.ConsumerGroup,
		"auto.offset.reset":     r.KafkaOffset,
		"session.timeout.ms":    r.KafkaSession,
		"heartbeat.interval.ms": r.KafkaHeartbeat,
		"fetch.min.bytes":       r.KafkaFetch,
		"max.poll.interval.ms":  r.KafkaMaxpoll,
	})

	if err != nil {
		fmt.Printf("ERROR. Failed consuming the topics. Reason: %v;", err)
		os.Exit(1)
	}

	_ = read.SubscribeTopics([]string{r.Topic}, nil)

	return read
}

func (w KafkaConfig) NewConfluentWriter() *kafka.Producer {
	brokersJoined := strings.Join(w.Brokers, `,`)
	kafka_batch, _ := strconv.Atoi(w.KafkaBatch)
	kafka_linger, _ := strconv.Atoi(w.KafkaLinger)
	produce, _ := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":  brokersJoined,
		"sasl.mechanisms":    w.KafkaMechanisms,
		"security.protocol":  w.KafkaProtocol,
		"sasl.username":      w.KafkaUsr,
		"sasl.password":      w.KafkaPsw,
		"linger.ms":          kafka_linger,
		"batch.num.messages": kafka_batch,
	})
	return produce
}

func (w KafkaConfig) Send(ctx context.Context, key string, message interface{}) (context.Context, error) {
	span, _ := apm.StartSpan(ctx, "Send", "ProduceKafkaConfluent")
	defer span.End()

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	writer := KafkaConfig.NewConfluentWriter(KafkaConfig{
		Brokers:         w.Brokers,
		KafkaBatch:      w.KafkaBatch,
		KafkaLinger:     w.KafkaLinger,
		KafkaMechanisms: w.KafkaMechanisms,
		KafkaProtocol:   w.KafkaProtocol,
		KafkaUsr:        w.KafkaUsr,
		KafkaPsw:        w.KafkaPsw,
	})
	defer writer.Close()

	data, _ := json.Marshal(message)
	s := string(data)
	s = strings.ReplaceAll(s, `\`, ``)

	// CONFLUENT INC
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &w.Topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: data,
	}
	err := writer.Produce(msg, nil)

	tr := &trace.TraceHttp{
		Url:     strings.Join(w.Brokers, `,`),
		Request: message,
	}

	currentTrace = append(currentTrace, tr)

	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)
	if err != nil {
		return ctx, err
	}

	flush, _ := strconv.Atoi(w.KafkaFlush)
	writer.Flush(1 * flush)
	return ctx, nil
}
