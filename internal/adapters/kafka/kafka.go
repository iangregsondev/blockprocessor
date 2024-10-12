package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/iangregsondev/deblockprocessor/pkg/logger"
	"github.com/segmentio/kafka-go"
)

// Define a struct to represent the message
type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

// Client is the struct for handling Client connections and operations
type Client struct {
	logger        logger.Logger
	brokerAddress string

	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader
	mutex   sync.Mutex
}

// NewClient initializes the Client with the broker address
func NewClient(logger logger.Logger, brokerAddress string) Kafka {
	return &Client{
		logger:        logger,
		brokerAddress: brokerAddress,
		writers:       make(map[string]*kafka.Writer),
		readers:       make(map[string]*kafka.Reader),
	}
}

// PublishMessage sends a message to a specified Client topic
func (c *Client) PublishMessage(ctx context.Context, topic string, key, value []byte) error {
	c.mutex.Lock()

	writer, exists := c.writers[topic]
	if !exists {
		// Create a new writer for the topic if it doesn't exist
		writer = &kafka.Writer{
			Addr:     kafka.TCP(c.brokerAddress),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
		c.writers[topic] = writer
	}

	c.mutex.Unlock()

	err := writer.WriteMessages(
		ctx, kafka.Message{
			Key:   key,
			Value: value,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	c.logger.Debug(fmt.Sprintf("message sent successfully to topic %s", topic))

	return nil
}

// Subscribe subscribes to a specific Client topic
func (c *Client) Subscribe(ctx context.Context, topic string, groupID string, processMessage func(topic string, key, value []byte)) error {
	c.mutex.Lock()

	reader, exists := c.readers[topic]
	if !exists {
		// TODO: These are hardcoded values, should be configurable
		const (
			minBytes = 10e3 // 10KB
			maxBytes = 10e6 // 10MB
		)
		// Create a new reader for the topic if it doesn't exist
		reader = kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{c.brokerAddress},
				Topic:   topic,
				// TODO choose either the groupID or Partition!
				GroupID: groupID,
				// Partition: 0,
				MinBytes: minBytes,
				MaxBytes: maxBytes,
			},
		)
		c.readers[topic] = reader
	}
	c.mutex.Unlock()

	errCh := make(chan error, 1)

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					errCh <- nil

					return
				}

				errCh <- fmt.Errorf("failed to read message from topic %s: %w", topic, err)

				return
			}

			// Log the offset, partition, and other message details
			log.Printf("Consumed message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)

			processMessage(msg.Topic, msg.Key, msg.Value)
		}
	}()

	return <-errCh
}

// Close closes all Client connections
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, writer := range c.writers {
		if err := writer.Close(); err != nil {
			log.Printf("failed to close writer: %v", err)
		}
	}

	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			log.Printf("failed to close reader: %v", err)
		}
	}
}
