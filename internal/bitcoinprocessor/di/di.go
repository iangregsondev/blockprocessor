package di

import (
	"fmt"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/adapters/sqlite"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/app"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/config"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/repository/block"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/services/blockchain"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/services/database"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/blockdaemon"
	"github.com/iangregsondev/deblockprocessor/pkg/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/rpc"
)

func NewDI(logger logger.Logger, mask *masker.MaskerMarshaler, osWrapper oswrapper.OS, cfg *config.Config) (*app.App, error) {
	logger.Info("Loading and validating configuration...")

	loadedConfig, err := cfg.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	// mask the configuration because it contains sensitive information
	maskedLoadedConfig, err := mask.Struct(loadedConfig)
	if err != nil {
		return nil, fmt.Errorf("error masking configuration: %w", err)
	}

	logger.Info("Configuration loaded and validated successfully", "config", maskedLoadedConfig)

	// instantiate adapters
	dbAdapter := sqlite.NewSqliteDatabase(osWrapper, loadedConfig.Database.File)
	rpcClient := rpc.NewRPC(loadedConfig.Connection.RPCURL, loadedConfig.Connection.APIKey)
	kafkaAdapter := kafka.NewClient(loadedConfig.Kafka.Broker)

	// provider
	bcProvider := blockdaemon.NewProvider(rpcClient)

	// repositories
	blockRepository := block.NewBlockRepository(dbAdapter)

	// instantiate services
	databaseService := database.NewService(dbAdapter)

	blockchainService := blockchain.NewService(
		logger, blockRepository, kafkaAdapter, bcProvider,
		loadedConfig.Blockchain.Height, loadedConfig.Polling.HeightIntervalMilliseconds, loadedConfig.Polling.BlockIntervalMilliseconds,
		loadedConfig.Connection.RPCURL, loadedConfig.Connection.APIKey, loadedConfig.Kafka.Topic,
	)

	// return the application
	return app.NewApp(logger, mask, osWrapper, loadedConfig, blockchainService, databaseService), nil
}
