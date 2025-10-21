package segmentio

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/xanderxampp-be/franco/contextwrap"
	"github.com/xanderxampp-be/franco/trace"

	"github.com/segmentio/kafka-go"
	"go.elastic.co/apm"
)

type Ksegmentio struct {
	Brokers []string
	Topic   string
	Group   string
	writer  *kafka.Writer
}

func NewProducer(brokers []string, topic, group string) *Ksegmentio {
	kafkaconf := kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.CRC32Balancer{},
	}
	writer := kafka.NewWriter(kafkaconf)

	return &Ksegmentio{
		writer: writer,
	}
}

func (k *Ksegmentio) Produce(ctx context.Context, key string, message interface{}) (context.Context, error) {
	span, _ := apm.StartSpan(ctx, "Send", "ProduceMessageSegmentio")
	defer span.End()

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	data, _ := json.Marshal(message)
	s := string(data)
	s = strings.ReplaceAll(s, `\`, ``)

	// SEGMENTIO
	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}

	err := k.writer.WriteMessages(context.Background(), msg)

	tr := &trace.TraceHttp{
		Url:     strings.Join(k.Brokers, `,`),
		Request: message,
	}

	currentTrace = append(currentTrace, tr)

	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}
