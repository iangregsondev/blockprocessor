package config

type Config struct {
	Blockchain BlockchainConfig `mapstructure:"blockchain" validate:"required" mask:"struct"`
	Connection ConnectionConfig `mapstructure:"connection" validate:"required" mask:"struct"`
	Database   DatabaseConfig   `mapstructure:"database" validate:"required" mask:"struct"`
	Polling    PollingConfig    `mapstructure:"polling" validate:"required" mask:"struct"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required" mask:"struct"`
}

type BlockchainConfig struct {
	Height int `mapstructure:"height" validate:"required"`
}

type ConnectionConfig struct {
	RPCURL string `mapstructure:"rpc_url" validate:"required"`
	APIKey string `mapstructure:"api_key" validate:"required" mask:"password"`
}

type DatabaseConfig struct {
	File string `mapstructure:"file" validate:"required"`
}

type PollingConfig struct {
	BlockIntervalMilliseconds  int `mapstructure:"block_interval_milliseconds" validate:"required"`
	HeightIntervalMilliseconds int `mapstructure:"height_interval_milliseconds" validate:"required"`
}

type KafkaConfig struct {
	Broker string `mapstructure:"broker" validate:"required"`
	Topic  string `mapstructure:"topic" validate:"required"`
}
