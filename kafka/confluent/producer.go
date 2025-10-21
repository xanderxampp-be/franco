package confluent

import (
	"context"
	"encoding/json"

	"franco/contextwrap"
	"franco/log"
	"franco/trace"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.elastic.co/apm"
)

type Producer struct {
	producer *kafka.Producer
	Config   *kafka.ConfigMap
}

func NewProducer(brokersJoined, saslMechanisms, securityProtocol, saslUsername, saslPassword string, lingerMs, batchNumMessages int) *Producer {
	producerCfg := &kafka.ConfigMap{
		"bootstrap.servers":  brokersJoined,
		"linger.ms":          lingerMs,
		"batch.num.messages": batchNumMessages,
		"compression.type":   "gzip",
	}

	if saslMechanisms != "" && saslUsername != "" && saslPassword != "" && securityProtocol != "" {
		producerCfg.SetKey("sasl.mechanisms", saslMechanisms)
		producerCfg.SetKey("security.protocol", securityProtocol)
		producerCfg.SetKey("sasl.username", saslUsername)
		producerCfg.SetKey("sasl.password", saslPassword)
	}

	producer, err := kafka.NewProducer(producerCfg)
	if err != nil {
		log.LogDebug("Error creating producer: " + err.Error())
		panic(err)
	}

	Producer := &Producer{
		producer: producer,
		Config:   producerCfg,
	}
	return Producer

}

func (b *Producer) Produce(ctx context.Context, key, topic string, message interface{}, flush int) (context.Context, error) {
	span, _ := apm.StartSpan(ctx, "Send", "ProduceKafkaConfluent")
	defer span.End()

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	data, err := json.Marshal(message)
	if err != nil {
		log.LogDebug("Error marshalling message: " + err.Error())
		return ctx, err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
		Key:   []byte(key),
	}

	report := make(chan kafka.Event, 1) // buffer to ensure channel isn't blocked
	defer close(report)

	err = b.producer.Produce(msg, report)
	if err != nil {
		log.LogDebug("Error during produce: " + err.Error())
		return ctx, err
	}

	event := <-report
	m := event.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.LogDebug("Delivery failed: " + m.TopicPartition.Error.Error())
		return ctx, m.TopicPartition.Error
	}

	bootstrapServersRaw, err := b.Config.Get("bootstrap.servers", nil)
	if err != nil {
		log.LogDebug("Error getting bootstrap servers config: " + err.Error())
	}

	bootstrapServers, ok := bootstrapServersRaw.(string)
	if !ok {
		log.LogDebug("Error asserting bootstrap servers config")
	}

	tr := &trace.TraceHttp{
		Url:     bootstrapServers,
		Request: message,
	}

	currentTrace = append(currentTrace, tr)
	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)

	b.producer.Flush(flush)
	return ctx, nil
}

func (b *Producer) ProduceWithSchema(ctx context.Context, key, topic string, message []byte, flush int) (context.Context, error) {
	span, _ := apm.StartSpan(ctx, "Send", "ProduceKafkaConfluent")
	defer span.End()

	currentTrace := contextwrap.GetTraceFromContext(ctx)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
		Key:   []byte(key),
	}

	report := make(chan kafka.Event, 1) // buffer to ensure channel isn't blocked
	defer close(report)

	err := b.producer.Produce(msg, report)
	if err != nil {
		log.LogDebug("why error produce : " + err.Error())
		return ctx, err
	}

	event := <-report
	m := event.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		log.LogDebug("Delivery failed: " + m.TopicPartition.Error.Error())
		return ctx, m.TopicPartition.Error
	}

	bootstrapServersRaw, err := b.Config.Get("bootstrap.servers", nil)
	if err != nil {
		log.LogDebug("error get config " + err.Error())
	}

	bootstrapServers, ok := bootstrapServersRaw.(string)
	if !ok {
		log.LogDebug("error assertion bootstrap server trace")
	}

	tr := &trace.TraceHttp{
		Url:     bootstrapServers,
		Request: message,
	}

	currentTrace = append(currentTrace, tr)
	ctx = context.WithValue(ctx, contextwrap.TraceKey, currentTrace)

	b.producer.Flush(flush)
	return ctx, nil
}
