package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoin/blockprocessor/repository/block"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin"
)

type Service struct {
	logger       logger.Logger
	kafkaAdapter kafka.Kafka
	bcProvider   bitcoin.Provider
	repository   *block.Repository

	defaultHeight         int
	blockPollingInterval  int
	heightPollingInterval int
	longestChain          *int

	rpcURL     string
	apiKey     string
	kafkaTopic string

	height int
}

func NewService(
	logger logger.Logger,
	repository *block.Repository,
	kafkaAdapter kafka.Kafka,
	bcProvider bitcoin.Provider,
	defaultHeight int, heightPollingInterval int, blockPollingInterval int, rpcURL string, apiKey string, kafkaTopic string,
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

func (b *Service) StartHeightWorker(ctx context.Context, wg *sync.WaitGroup, heightCh chan int) error {
	b.logger.Info("Starting height worker...")

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
				b.logger.Info("Polling blockchain for block count...")

				blockCount, err := b.bcProvider.GetBlockCount(ctx)
				if err != nil {
					b.logger.Error("Error getting block count", "error", err)

					break
				}

				b.logger.Debug("Got block count", "count", blockCount.Result)

				b.longestChain = &blockCount.Result

			case newHeight := <-heightCh:
				// Update the height in the database
				err := b.repository.CreateOrUpdateBlockHeight(newHeight)
				if err != nil {
					b.logger.Error("Error creating/updating block height in database", "error", err)
				}

			case <-ctx.Done():
				b.logger.Info("Stopping height worker due to context cancellation...")

				return
			}
		}
	}()

	return nil
}

func (b *Service) StartBlockWorker(ctx context.Context, wg *sync.WaitGroup, heightCh chan int) error {
	b.logger.Info("Starting block worker...")

	// Obtain height
	height, err := b.repository.GetBlockHeight()
	if err != nil {
		if errors.Is(err, block.ErrBlockHeightNotFound) {
			// Handle the case where no BlockHeight was found
			err := b.repository.CreateOrUpdateBlockHeight(b.defaultHeight)
			if err != nil {
				return fmt.Errorf("error creating/updating block height in database: %w", err)
			}

			// Set the height to the default height which is now in the DB
			b.height = b.defaultHeight
		} else {
			return fmt.Errorf("error getting block height from database: %w", err)
		}
	} else {
		b.height = height.Height
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
				err := b.pollBlockchain(ctx, heightCh)
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

func (b *Service) pollBlockchain(ctx context.Context, heightCh chan int) error {
	b.logger.Info("Polling blockchain...")

	if b.longestChain == nil || b.height >= *b.longestChain {
		if b.longestChain == nil {
			b.logger.Debug("No idea what the longest chain is, awaiting longest chain :-(")
		} else {
			b.logger.Info("Height is at or beyond the longest chain, waiting for new blocks...")
		}

		return nil
	}

	blockHash, err := b.bcProvider.GetBlockHash(ctx, b.height)
	if err != nil {
		return fmt.Errorf("error getting block hash: %w", err)
	}

	if blockHash.Error != nil {
		return fmt.Errorf("error returned from getting block hash: %s", blockHash.Error)
	}

	header, err := b.bcProvider.GetBlockHeader(ctx, blockHash.Result)
	if err != nil {
		return fmt.Errorf("error getting block header: %w", err)
	}

	if header.Error != nil {
		return fmt.Errorf("error returned from getting block header: %s", header.Error)
	}

	value, err := json.Marshal(header.Result)
	if err != nil {
		return fmt.Errorf("error marshalling block header: %w", err)
	}

	b.logger.Info("Publishing block to Kafka", "topic", b.kafkaTopic, "height", b.height, "hash", header.Result.Hash)

	err = b.kafkaAdapter.PublishMessage(ctx, b.kafkaTopic, []byte(header.Result.Hash), value, kafka.PublishOptions{})
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	// Increment the blockchain height after successful publishing
	b.height++
	b.logger.Debug("Incremented blockchain height", "height", b.height)

	// Send the new height to the height worker via channel to ensure its not blocking!
	heightCh <- b.height

	return nil
}
