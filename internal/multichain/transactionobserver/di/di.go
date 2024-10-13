package di

import (
	"fmt"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/adapters/rpc"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/app"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/config"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/helper"
	// configmodel "github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/processors"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/processors/bitcoin"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/processors/ethereum"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/processors/solana"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/transaction"
	iowrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/io"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	bitcoinblockdaemon "github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/blockdaemon"
	ethereumblockdaemon "github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/blockdaemon"
	solanablockdaemon "github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana/blockdaemon"
)

func NewDI(logger logger.Logger, mask *masker.MaskerMarshaler, ioWrapper iowrapper.IO, osWrapper oswrapper.OS, cfg *config.Config) (*app.App, error) {
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

	level, err := logger.ParseLogLevel(loadedConfig.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("error parsing log level: %w", err)
	}

	logger.SetLogLevel(level)

	// instantiate adapters

	// instantiate rpc clients, we only support a fixed amount of rpc clients for now
	rpcBitcoinClient := rpc.NewRPC(
		logger, ioWrapper, loadedConfig.Connection.Bitcoin.RPCURL, loadedConfig.Connection.Bitcoin.APIKey, loadedConfig.HTTP.MaxRetryOnError,
		loadedConfig.HTTP.RetryDelayMilliseconds, loadedConfig.HTTP.ReportRetryAttempts,
	)

	rpcEthereumClient := rpc.NewRPC(
		logger, ioWrapper, loadedConfig.Connection.Ethereum.RPCURL, loadedConfig.Connection.Ethereum.APIKey, loadedConfig.HTTP.MaxRetryOnError,
		loadedConfig.HTTP.RetryDelayMilliseconds, loadedConfig.HTTP.ReportRetryAttempts,
	)

	rpcSolanaClient := rpc.NewRPC(
		logger, ioWrapper, loadedConfig.Connection.Solana.RPCURL, loadedConfig.Connection.Solana.APIKey, loadedConfig.HTTP.MaxRetryOnError,
		loadedConfig.HTTP.RetryDelayMilliseconds, loadedConfig.HTTP.ReportRetryAttempts,
	)

	kafkaAdapter := kafka.NewClient(logger, loadedConfig.Kafka.Broker)

	// Build available providers, we support a fixed amount of providers for now
	bcBitcoinProvider := bitcoinblockdaemon.NewProvider(rpcBitcoinClient)
	bcEthereumProvider := ethereumblockdaemon.NewProvider(rpcEthereumClient)
	bcSolanaProvider := solanablockdaemon.NewProvider(rpcSolanaClient)

	// Build provider map
	providerMap := make(map[string]any)
	providerMap["bitcoin"] = bcBitcoinProvider
	providerMap["ethereum"] = bcEthereumProvider
	providerMap["solana"] = bcSolanaProvider

	processorMap := make(map[string]processors.Processor)

	// Build available processors, we only support a fixed amount of processors for now
	for _, worker := range loadedConfig.Worker {
		chain := worker.Chain

		if _, exists := providerMap[chain]; !exists {
			return nil, fmt.Errorf("no provider support for chain: %s", chain)
		}

		users := helper.GetUsersByChain(loadedConfig.Users, chain)

		var processor processors.Processor

		switch chain {
		case "bitcoin":
			processor = bitcoin.NewProcessor(logger, bcBitcoinProvider, users)
		case "ethereum":
			processor = ethereum.NewProcessor(logger, bcEthereumProvider, users)
		case "solana":
			processor = solana.NewProcessor(logger, bcSolanaProvider, users)
		default:
			return nil, fmt.Errorf("unsupported chain, no processors are available: %s", chain)
		}

		processorMap[chain] = processor
	}

	// instantiate services
	transactionService := transaction.NewService(
		logger, kafkaAdapter, loadedConfig.Kafka.ChainTopics, loadedConfig.Worker, loadedConfig.Users, processorMap,
	)

	// return the application
	return app.NewApp(logger, mask, osWrapper, loadedConfig, transactionService), nil
}
