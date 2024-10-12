package config

type Config struct {
	Connection ConnectionConfig `mapstructure:"connection" validate:"required" mask:"struct"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required" mask:"struct"`
	Worker     WorkerConfig     `mapstructure:"worker" validate:"required" mask:"struct"`
}

type ConnectionConfig struct {
	RPCURL string `mapstructure:"rpc_url" validate:"required"`
	APIKey string `mapstructure:"api_key" validate:"required" mask:"password"`
}

type KafkaConfig struct {
	Broker             string `mapstructure:"broker" validate:"required"`
	BlockTopic         string `mapstructure:"block_topic" validate:"required"`
	BlockConsumerGroup string `mapstructure:"block_consumer_group" validate:"required"`
	TransactionTopic   string `mapstructure:"transaction_topic" validate:"required"`
}

type WorkerConfig struct {
	Total int `mapstructure:"total" validate:"required"`
}
