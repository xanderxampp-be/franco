package confluent

import (
	"context"
	"fmt"

	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/sync/semaphore"
)

type Consumer struct {
	Consumer *kafka.Consumer
	Config   *kafka.ConfigMap
}

func NewConsumer(brokersJoined, groupId, autoOffsetReset, sessionTimeoutMs, heartbeatIntervalMs, fetchMinBytes, maxPollIntervalMs, saslMechanisms, securityProtocol, saslUsername, saslPassword string) *Consumer {
	consumerCfg := &kafka.ConfigMap{
		"bootstrap.servers":     brokersJoined,
		"group.id":              groupId,
		"auto.offset.reset":     autoOffsetReset,
		"session.timeout.ms":    sessionTimeoutMs,
		"heartbeat.interval.ms": heartbeatIntervalMs,
		"fetch.min.bytes":       fetchMinBytes,
		"max.poll.interval.ms":  maxPollIntervalMs,
	}

	if saslMechanisms != "" && saslUsername != "" && saslPassword != "" && securityProtocol != "" {
		consumerCfg.SetKey("sasl.mechanisms", saslMechanisms)
		consumerCfg.SetKey("security.protocol", securityProtocol)
		consumerCfg.SetKey("sasl.username", saslUsername)
		consumerCfg.SetKey("sasl.password", saslPassword)
	}

	consumer, err := kafka.NewConsumer(consumerCfg)
	if err != nil {
		fmt.Println("Failed to create Kafka consumer:  ", err.Error())
	}

	return &Consumer{
		Consumer: consumer,
		Config:   consumerCfg,
	}
}

// HandleConsume handle consumption with semaphore and slowdown if semaphore out of room.
// The purpose is to achieve simple consumption method with limitation related to semaphore size.
// fn as sub-function, is the main logic that can be inserted in this method to customize the outcome of consumption
func (b *Consumer) HandleConsume(sourceTopic string, semaphoreSize int, producer *Producer, slowDownMilisec int, fn func(ctx context.Context, msgVal []byte, topic string, offset string, partition string)) {
	err := b.Consumer.Subscribe(sourceTopic, nil)
	if err != nil {
		fmt.Println("Failed to subscribe to topic: ", err.Error())
	}

	sem := semaphore.NewWeighted(int64(semaphoreSize))
	delayedMessages := make(chan *kafka.Message, 100)

	// Goroutine to handle delayed messages
	go func() {
		for msg := range delayedMessages {
			sem.Acquire(context.Background(), 1)
			go func(msg *kafka.Message) {
				defer sem.Release(1)
				topicName := *msg.TopicPartition.Topic
				topicOffset := msg.TopicPartition.Offset.String()
				topicPartition := strconv.Itoa(int(msg.TopicPartition.Partition))
				fn(context.Background(), msg.Value, topicName, topicOffset, topicPartition)
			}(msg)
		}
	}()

	// start consuming from topic
	for {
		ctx := context.Background()

		msg, err := b.Consumer.ReadMessage(-1)
		if err != nil {
			fmt.Println("error on polling read message : ", err.Error())
			continue
		}

		// acquire semaphore, if semaphore has not enough room, slowdoon the loop and reproduce message back to source topic
		if ok := sem.TryAcquire(1); !ok {
			fmt.Println("max semaphore detected")
			time.Sleep(time.Duration(slowDownMilisec) * time.Millisecond)
			delayedMessages <- msg
			continue
		}

		go func(msg *kafka.Message) {
			// release semaphore
			defer sem.Release(1)
			topicName := *msg.TopicPartition.Topic
			topicOffset := msg.TopicPartition.Offset.String()
			topicPartition := strconv.Itoa(int(msg.TopicPartition.Partition))
			fn(ctx, msg.Value, topicName, topicOffset, topicPartition)
		}(msg)
	}
}
