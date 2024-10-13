package config

type Config struct {
	Connection ConnectionConfig `mapstructure:"connection" validate:"required" mask:"struct"`
	HTTP       HTTPConfig       `mapstructure:"http" validate:"required" mask:"struct"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required" mask:"struct"`
	Logger     LoggerConfig     `mapstructure:"logger" validate:"required" mask:"struct"`
	Worker     []WorkerConfig   `mapstructure:"workers" validate:"required" mask:"struct"`
	Users      []UserConfig     `mapstructure:"users" validate:"required" mask:"struct"`
}

type ConnectionConfig struct {
	Bitcoin  ConnectionChainConfig `mapstructure:"bitcoin" validate:"required" mask:"struct"`
	Ethereum ConnectionChainConfig `mapstructure:"ethereum" validate:"required" mask:"struct"`
	Solana   ConnectionChainConfig `mapstructure:"solana" validate:"required" mask:"struct"`
}

type ConnectionChainConfig struct {
	RPCURL string `mapstructure:"rpc_url" validate:"required"`
	APIKey string `mapstructure:"api_key" validate:"required" mask:"password"`
}

type HTTPConfig struct {
	MaxRetryOnError        int  `mapstructure:"max_retry_on_error" validate:"required"`
	RetryDelayMilliseconds int  `mapstructure:"retry_delay_milliseconds" validate:"required"`
	ReportRetryAttempts    bool `mapstructure:"report_retry_attempts"`
}

type KafkaConfig struct {
	Broker      string             `mapstructure:"broker" validate:"required"`
	ChainTopics []ChainTopicConfig `mapstructure:"chain_topics" validate:"required"`
}

type ChainTopicConfig struct {
	Chain         string  `mapstructure:"chain" validate:"required"`
	Topic         string  `mapstructure:"topic" validate:"required"`
	ConsumerGroup *string `mapstructure:"consumer_group"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level" validate:"required" `
}

type WorkerConfig struct {
	Chain string `mapstructure:"chain" validate:"required"`
	Total int    `mapstructure:"total" validate:"required"`
}

type UserConfig struct {
	Name    string `mapstructure:"name" validate:"required"`
	Chain   string `mapstructure:"chain" validate:"required"`
	Address string `mapstructure:"address" validate:"required"`
}
