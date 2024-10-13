package app

import (
	"context"
	"sync"
	"syscall"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/common/signal"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/transaction"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
)

type App struct {
	logger logger.Logger
	config *config.Config

	transactionService *transaction.Service

	osWrapper oswrapper.OS

	mask *masker.MaskerMarshaler
}

func NewApp(
	logger logger.Logger, mask *masker.MaskerMarshaler, osWrapper oswrapper.OS, config *config.Config, transactionService *transaction.Service,
) *App {
	return &App{
		logger:             logger,
		config:             config,
		transactionService: transactionService,
		osWrapper:          osWrapper,
		mask:               mask,
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

	// Start services
	a.transactionService.StartTransactionQueueWorker(ctx, &wg)

	wg.Wait()

	return nil
}
