package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/ethereum/blockprocessor/repository/block"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum"
)

type Service struct {
	logger       logger.Logger
	kafkaAdapter kafka.Kafka
	bcProvider   ethereum.Provider
	repository   *block.Repository

	defaultHeight         int64
	blockPollingInterval  int
	heightPollingInterval int
	latestBlockNumber     *int64

	rpcURL     string
	apiKey     string
	kafkaTopic string

	currentBlockNumber int64
}

func NewService(
	logger logger.Logger,
	repository *block.Repository,
	kafkaAdapter kafka.Kafka,
	bcProvider ethereum.Provider,
	defaultHeight int64, heightPollingInterval int, blockPollingInterval int, rpcURL string, apiKey string, kafkaTopic string,
) *Service {
	return &Service{
		logger:                logger,
		repository:            repository,
		kafkaAdapter:          kafkaAdapter,
		bcProvider:            bcProvider,
		defaultHeight:         defaultHeight,
		blockPollingInterval:  blockPollingInterval,
		heightPollingInterval: heightPollingInterval,
		rpcURL:                rpcURL,
		apiKey:                apiKey,
		kafkaTopic:            kafkaTopic,
	}
}

func (b *Service) StartCurrentBlockNumberWorker(ctx context.Context, wg *sync.WaitGroup, blockNumberCh chan int64) error {
	b.logger.Info("Starting currentBlockNumber (block number) worker...")

	wg.Add(1)

	go func() {
		defer wg.Done()

		// Create a ticker for polling interval
		ticker := time.NewTicker(time.Duration(b.blockPollingInterval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Perform the polling work here
				b.logger.Info("Polling blockchain for latest block number...")

				blockNumber, err := b.bcProvider.GetBlockNumber(ctx)
				if err != nil {
					b.logger.Error("Error getting block count", "error", err)

					break
				}

				latestBlockNumber, err := b.convertHexToDecimal(blockNumber.Result)
				if err != nil {
					b.logger.Error("error converting hex to decimal", "hex", "blockNumber.Result", "error", err)

					break
				}

				b.logger.Debug("Got block number", "hex", blockNumber.Result, "number", *latestBlockNumber)

				b.latestBlockNumber = latestBlockNumber

			case newBlockNumber := <-blockNumberCh:
				// Update the currentBlockNumber in the database
				err := b.repository.CreateOrUpdateBlockNumber(newBlockNumber)
				if err != nil {
					b.logger.Error("Error creating/updating block currentBlockNumber in database", "error", err)
				}

			case <-ctx.Done():
				b.logger.Info("Stopping currentBlockNumber worker due to context cancellation...")

				return
			}
		}
	}()

	return nil
}

func (b *Service) StartBlockWorker(ctx context.Context, wg *sync.WaitGroup, blockNumberCh chan int64) error {
	b.logger.Info("Starting block worker...")

	// Obtain currentBlockNumber
	blockNumber, err := b.repository.GetLatestBlockNumber()
	if err != nil {
		if errors.Is(err, block.ErrBlockNumberNotFound) {
			// Handle the case where no BlockHeight was found
			err := b.repository.CreateOrUpdateBlockNumber(b.defaultHeight)
			if err != nil {
				return fmt.Errorf("error creating/updating block currentBlockNumber in database: %w", err)
			}

			// Set the currentBlockNumber to the default currentBlockNumber which is now in the DB
			b.currentBlockNumber = b.defaultHeight
		} else {
			return fmt.Errorf("error getting block currentBlockNumber from database: %w", err)
		}
	} else {
		b.currentBlockNumber = blockNumber.BlockNumber
	}

	// Polling loop
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Create a ticker for polling interval
		ticker := time.NewTicker(time.Duration(b.blockPollingInterval) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Perform the polling work here
				err := b.pollBlockchain(ctx, blockNumberCh)
				if err != nil {
					b.logger.Error("Error polling blockchain", "error", err)
				}

			case <-ctx.Done():
				b.logger.Info("Stopping blockchain worker due to context cancellation...")

				return
			}
		}
	}()

	return nil
}

func (b *Service) pollBlockchain(ctx context.Context, blockNumberCh chan int64) error {
	b.logger.Info("Polling blockchain...")

	if b.latestBlockNumber == nil || b.currentBlockNumber >= *b.latestBlockNumber {
		if b.latestBlockNumber == nil {
			b.logger.Debug("No idea what the longest chain is, awaiting longest chain :-(")
		} else {
			b.logger.Info("BlockNumber is at or beyond the longest chain, waiting for new blocks...")
		}

		return nil
	}

	currentBlockNumberHex := fmt.Sprintf("0x%x", b.currentBlockNumber)

	blk, err := b.bcProvider.GetBlockByNumber(ctx, currentBlockNumberHex, false)
	if err != nil {
		return fmt.Errorf("error getting block by number: %w", err)
	}

	if blk.Error != nil {
		return fmt.Errorf("error returned from getting block by number: %v", blk.Error)
	}

	value, err := json.Marshal(blk.Result)
	if err != nil {
		return fmt.Errorf("error marshalling block: %w", err)
	}

	b.logger.Info("Publishing block to Kafka", "topic", b.kafkaTopic, "height", b.currentBlockNumber, "hash", blk.Result.Hash)

	err = b.kafkaAdapter.PublishMessage(ctx, b.kafkaTopic, []byte(blk.Result.Hash), value)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	// Increment the blockchain currentBlockNumber after successful publishing
	b.currentBlockNumber++
	b.logger.Debug("Incremented blockchain currentBlockNumber", "currentBlockNumber", b.currentBlockNumber)

	// Send the new currentBlockNumber to the currentBlockNumber worker via channel to ensure it's not blocking!
	blockNumberCh <- b.currentBlockNumber

	return nil
}

func (b *Service) convertHexToDecimal(hex string) (*int64, error) {
	decimal, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		b.logger.Error("Error converting from hex to decimal", "hex", hex, "error", err)

		return nil, err
	}

	return &decimal, nil
}
