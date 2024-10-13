package config

type Config struct {
	Connection ConnectionConfig `mapstructure:"connection" validate:"required" mask:"struct"`
	HTTP       HTTPConfig       `mapstructure:"http" validate:"required" mask:"struct"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required" mask:"struct"`
	Logger     LoggerConfig     `mapstructure:"logger" validate:"required" mask:"struct"`
	Worker     WorkerConfig     `mapstructure:"worker" validate:"required" mask:"struct"`
}

type ConnectionConfig struct {
	RPCURL string `mapstructure:"rpc_url" validate:"required"`
	APIKey string `mapstructure:"api_key" validate:"required" mask:"password"`
}

type HTTPConfig struct {
	MaxRetryOnError        int  `mapstructure:"max_retry_on_error" validate:"required"`
	RetryDelayMilliseconds int  `mapstructure:"retry_delay_milliseconds" validate:"required"`
	ReportRetryAttempts    bool `mapstructure:"report_retry_attempts"`
}

type KafkaConfig struct {
	Broker             string  `mapstructure:"broker" validate:"required"`
	BlockTopic         string  `mapstructure:"block_topic" validate:"required"`
	BlockConsumerGroup *string `mapstructure:"block_consumer_group"`
	TransactionTopic   string  `mapstructure:"transaction_topic" validate:"required"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level" validate:"required" `
}

type WorkerConfig struct {
	Total int `mapstructure:"total" validate:"required"`
}
