package config

type Config struct {
	Blockchain BlockchainConfig `mapstructure:"blockchain" validate:"required" mask:"struct"`
	Connection ConnectionConfig `mapstructure:"connection" validate:"required" mask:"struct"`
	Database   DatabaseConfig   `mapstructure:"database" validate:"required" mask:"struct"`
	HTTP       HTTPConfig       `mapstructure:"http" validate:"required" mask:"struct"`
	Kafka      KafkaConfig      `mapstructure:"kafka" validate:"required" mask:"struct"`
	Logger     LoggerConfig     `mapstructure:"logger" validate:"required" mask:"struct"`
	Polling    PollingConfig    `mapstructure:"polling" validate:"required" mask:"struct"`
}

type BlockchainConfig struct {
	SeedBlockNumber int `mapstructure:"seed_block_number" validate:"required"`
}

type ConnectionConfig struct {
	RPCURL string `mapstructure:"rpc_url" validate:"required"`
	APIKey string `mapstructure:"api_key" validate:"required" mask:"password"`
}

type DatabaseConfig struct {
	File string `mapstructure:"file" validate:"required"`
}

type HTTPConfig struct {
	MaxRetryOnError        int  `mapstructure:"max_retry_on_error" validate:"required"`
	RetryDelayMilliseconds int  `mapstructure:"retry_delay_milliseconds" validate:"required"`
	ReportRetryAttempts    bool `mapstructure:"report_retry_attempts"`
}

type KafkaConfig struct {
	Broker string `mapstructure:"broker" validate:"required"`
	Topic  string `mapstructure:"topic" validate:"required"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level" validate:"required" `
}

type PollingConfig struct {
	BlockIntervalMilliseconds  int `mapstructure:"block_interval_milliseconds" validate:"required"`
	HeightIntervalMilliseconds int `mapstructure:"height_interval_milliseconds" validate:"required"`
}
