package kafka

import "context"

type Kafka interface {
	PublishMessage(ctx context.Context, topic string, key, value []byte, options PublishOptions) error
	Subscribe(ctx context.Context, topic string, processMessage func(topic string, key, value []byte), options SubscriptionOptions) error
}
