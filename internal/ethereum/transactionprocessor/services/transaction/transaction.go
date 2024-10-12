package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/models/response"
)

type Service struct {
	logger       logger.Logger
	kafkaAdapter kafka.Kafka
	bcProvider   ethereum.Provider

	rpcURL                string
	apiKey                string
	kafkaBlockTopic       string
	kafkaBlockConsumer    string
	kafkaTransactionTopic string
}

func NewService(
	logger logger.Logger,
	kafkaAdapter kafka.Kafka,
	bcProvider ethereum.Provider,
	rpcURL string, apiKey string, kafkaBlockTopic string, kafkaBlockConsumer string, kafkaTransactionTopic string,
) *Service {
	return &Service{
		logger:                logger,
		kafkaAdapter:          kafkaAdapter,
		bcProvider:            bcProvider,
		rpcURL:                rpcURL,
		apiKey:                apiKey,
		kafkaBlockTopic:       kafkaBlockTopic,
		kafkaBlockConsumer:    kafkaBlockConsumer,
		kafkaTransactionTopic: kafkaTransactionTopic,
	}
}

func (b *Service) StartTransactionQueueWorker(ctx context.Context, wg *sync.WaitGroup, kafkaMessageCh chan kafka.Message) {
	b.logger.Info("Starting transaction queue worker...")

	defer wg.Done()

	wg.Add(1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("Stopping transaction queue worker due to context cancellation...")

				return
			default:
				err := b.kafkaAdapter.Subscribe(
					ctx, b.kafkaBlockTopic, b.kafkaBlockConsumer, func(topic string, key, value []byte) {
						// Process the message here
						b.logger.Info("Received message", "topic", topic, "key", key)
						kafkaMessageCh <- kafka.Message{Topic: topic, Key: key, Value: value}
					},
				)
				if err != nil {
					b.logger.Info("Error received from subscription to kafka topic", err)
				}
			}
		}
	}()
}

func (b *Service) StartWorkerPool(ctx context.Context, numWorkers int, wg *sync.WaitGroup, kafkaMessageCh chan kafka.Message) {
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go b.worker(ctx, i, wg, kafkaMessageCh)
	}
}

func (b *Service) worker(ctx context.Context, workerID int, wg *sync.WaitGroup, kafkaMessageCh chan kafka.Message) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			b.logger.Info(fmt.Sprintf("Worker %d: Context Cancelled, stopping...", workerID))

			return
		case msg, ok := <-kafkaMessageCh:
			if !ok {
				b.logger.Info(fmt.Sprintf("Worker %d: No more tasks, exiting...", workerID))

				return
			}

			b.logger.Info(fmt.Sprintf("Worker %d: Processing %s", workerID, msg))

			err := b.processTransactions(ctx, workerID, msg)
			if err != nil {
				b.logger.Error(fmt.Sprintf("Worker %d: error processing transactions, ", workerID), "error", err)
			}
		}
	}
}

func (b *Service) processTransactions(ctx context.Context, workerID int, msg kafka.Message) error {
	var blockByNumber response.BlockByNumber

	err := json.Unmarshal(msg.Value, &blockByNumber)
	if err != nil {
		return fmt.Errorf("worker %d: error marshalling block by number: %w", workerID, err)
	}

	block, err := b.bcProvider.GetBlockByHash(ctx, blockByNumber.Hash, true)
	if err != nil {
		return fmt.Errorf("worker %d: error getting block hash: %w", workerID, err)
	}

	if block.Error != nil {
		return fmt.Errorf("worker %d: error returned from getting block: %v", workerID, block.Error)
	}

	for _, tx := range block.Result.Transactions {
		rawTx, err := b.bcProvider.GetTransactionByHash(ctx, tx.Hash)
		if err != nil {
			return fmt.Errorf("worker %d: error getting block hash: %w", workerID, err)
		}

		if rawTx.Error != nil {
			return fmt.Errorf("worker %d: error returned from getting raw transaction: %v", workerID, rawTx.Error)
		}

		value, err := json.Marshal(rawTx.Result)
		if err != nil {
			return fmt.Errorf("worker %d: error marshalling raw transaction: %w", workerID, err)
		}

		err = b.kafkaAdapter.PublishMessage(ctx, b.kafkaTransactionTopic, []byte(tx.Hash), value)
		if err != nil {
			return fmt.Errorf("worker %d: error publishing message: %w", workerID, err)
		}
	}

	return nil
}