package app

import (
	"context"
	"sync"
	"syscall"

	"github.com/ggwhite/go-masker/v2"
	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/common/signal"
	"github.com/iangregsondev/deblockprocessor/internal/ethereum/transactionprocessor/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/ethereum/transactionprocessor/services/transaction"
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

	kafkaMessageCh := make(chan kafka.Message)

	// Start services
	a.transactionService.StartWorkerPool(ctx, a.config.Worker.Total, &wg, kafkaMessageCh)
	a.transactionService.StartTransactionQueueWorker(ctx, &wg, kafkaMessageCh)

	wg.Wait()

	return nil
}
