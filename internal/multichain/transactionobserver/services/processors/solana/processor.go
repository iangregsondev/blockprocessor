package solana

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana"
)

// const decimalPlaces = 9

type Processor struct {
	logger     logger.Logger
	bcProvider solana.Provider

	users []config.UserConfig
}

func NewProcessor(logger logger.Logger, bcProvider solana.Provider, users []config.UserConfig) *Processor {
	return &Processor{
		logger:     logger,
		bcProvider: bcProvider,
		users:      users,
	}
}

func (p *Processor) Process(ctx context.Context, workerID int, msg kafka.Message) error {
	// TODO: Implement Solana transaction processing
	_, _, _ = ctx, workerID, msg //nolint:dogsled

	p.logger.Debug("Processing transaction")

	return nil
}
