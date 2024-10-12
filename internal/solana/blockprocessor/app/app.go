package app

import (
	"context"
	"fmt"
	"sync"
	"syscall"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/common/signal"
	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/services/blockchain"
	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/services/database"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
)

type App struct {
	logger logger.Logger
	config *config.Config

	blockchainService *blockchain.Service
	databaseService   *database.Service

	osWrapper oswrapper.OS

	mask *masker.MaskerMarshaler
}

func NewApp(
	logger logger.Logger, mask *masker.MaskerMarshaler, osWrapper oswrapper.OS, config *config.Config, blockchainService *blockchain.Service,
	databaseService *database.Service,
) *App {
	return &App{
		logger:            logger,
		config:            config,
		blockchainService: blockchainService,
		databaseService:   databaseService,
		osWrapper:         osWrapper,
		mask:              mask,
	}
}

// Run starts the application's main logic.
func (a *App) Run() error {
	var wg sync.WaitGroup

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build and run the signal handler, passing the cancel function directly
	signalHandler := signal.BuildSignalHandler(a.logger, cancel, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	signalHandler(ctx)

	// Set up the database
	err := a.databaseService.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup the database: %w", err)
	}

	blockNumberCh := make(chan int64)

	// Start services
	err = a.blockchainService.StartBlockHeightWorker(ctx, &wg, blockNumberCh)
	if err != nil {
		return fmt.Errorf("failed to start height worker: %w", err)
	}

	err = a.blockchainService.StartBlockWorker(ctx, &wg, blockNumberCh)
	if err != nil {
		return fmt.Errorf("failed to start block worker: %w", err)
	}

	wg.Wait()

	return nil
}
