package kafka

import "context"

type Kafka interface {
	PublishMessage(ctx context.Context, topic string, key, value []byte) error
	Subscribe(ctx context.Context, topic string, processMessage func(topic string, key, value []byte)) error
}
