package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/segmentio/kafka-go"
)

// Message is the struct for handling Kafka messages
type Message struct {
	Topic string
	Key   []byte
	Value []byte
}

// SubscriptionOptions is the struct for handling subscription options
type SubscriptionOptions struct {
	GroupID     *string
	Partition   *int
	StartOffset *int64
}

// PublishOptions is the struct for handling publish options
type PublishOptions struct {
	Partition *int
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
func (c *Client) PublishMessage(ctx context.Context, topic string, key, value []byte, options PublishOptions) error {
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

	message := kafka.Message{
		Key:   key,
		Value: value,
	}

	// Set Partition if specified in the options
	if options.Partition != nil {
		message.Partition = *options.Partition
	}

	err := writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to topic %s: %w", topic, err)
	}

	c.logger.Debug(fmt.Sprintf("message sent successfully to topic %s", topic))

	return nil
}

// Subscribe subscribes to a specific Client topic
func (c *Client) Subscribe(ctx context.Context, topic string, processMessage func(topic string, key, value []byte), options SubscriptionOptions) error {
	c.mutex.Lock()

	reader, exists := c.readers[topic]
	if !exists {
		// TODO: These are hardcoded values, should be configurable
		const (
			minBytes = 10e3 // 10KB
			maxBytes = 10e6 // 10MB
		)

		readerConfig := kafka.ReaderConfig{
			Brokers:  []string{c.brokerAddress},
			Topic:    topic,
			MinBytes: minBytes,
			MaxBytes: maxBytes,
		}

		// Set GroupID if provided
		if options.GroupID != nil {
			readerConfig.GroupID = *options.GroupID
		}

		// Set Partition if provided, only if GroupID is not set
		if options.Partition != nil && options.GroupID == nil {
			readerConfig.Partition = *options.Partition
		}

		// Set StartOffset only if it's explicitly provided
		if options.StartOffset != nil && options.GroupID == nil {
			readerConfig.StartOffset = *options.StartOffset
		}

		reader = kafka.NewReader(readerConfig)
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
			c.logger.Debug(fmt.Sprintf("Consumed message from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset))

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
			c.logger.Error(fmt.Sprintf("failed to close writer: %v", err))
		}
	}

	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			c.logger.Error(fmt.Sprintf("failed to close reader: %v", err))
		}
	}
}
