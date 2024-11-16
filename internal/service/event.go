package eventservice

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

type Event struct {
	Operation string
	Entity    string
	Data      interface{}
}

type Producer interface {
	Produce(event *Event)
	Close() error
}

type KafkaProducer struct {
	kafkaWriter *kafka.Writer
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaProducer{
		kafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *KafkaProducer) Produce(event *Event) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		log.Printf("Produced Event: %+v\n", event)

		// Serialize the event
		eventData, err := json.Marshal(event)
		if err != nil {
			log.Printf("Failed to serialize event: %v\n", err)
			return
		}

		// Send the event to Kafka
		err = p.kafkaWriter.WriteMessages(p.ctx,
			kafka.Message{
				Value: eventData,
			},
		)
		if err != nil {
			log.Printf("Failed to send event to Kafka: %v\n", err)
		}
	}()
}

func (p *KafkaProducer) Close() error {
	p.cancel()
	p.wg.Wait()
	return p.kafkaWriter.Close()
}

type NoOpProducer struct{}

func (p *NoOpProducer) Produce(event *Event) {
	log.Printf("NoOpProducer: %+v\n", event)
}
func (p *NoOpProducer) Close() error { return nil }
