package transaction

import (
	"context"
	"fmt"
	"sync"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/services/processors"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
)

type Service struct {
	logger       logger.Logger
	kafkaAdapter kafka.Kafka

	chainTopics []config.ChainTopicConfig
	workers     []config.WorkerConfig
	users       []config.UserConfig

	chainChannelMap map[string]chan kafka.Message

	processorMap map[string]processors.Processor
}

func NewService(
	logger logger.Logger, kafkaAdapter kafka.Kafka, chainTopics []config.ChainTopicConfig, workers []config.WorkerConfig, users []config.UserConfig,
	processorMap map[string]processors.Processor,
) *Service {
	return &Service{
		logger:          logger,
		kafkaAdapter:    kafkaAdapter,
		chainTopics:     chainTopics,
		workers:         workers,
		users:           users,
		processorMap:    processorMap,
		chainChannelMap: make(map[string]chan kafka.Message, len(workers)),
	}
}

func (b *Service) StartTransactionQueueWorker(ctx context.Context, wg *sync.WaitGroup) {
	b.logger.Info("Starting transaction queue worker...")

	// TODO - Move this to a config file
	bufferedChannelSize := 200

	// Create a channel for each worker chain
	for _, worker := range b.workers {
		chain := worker.Chain
		if _, exists := b.chainChannelMap[chain]; !exists {
			// Create a new channel for the chain
			b.chainChannelMap[chain] = make(chan kafka.Message, bufferedChannelSize)
			b.logger.Info("Created channel for chain", "chain", chain)
		}
	}

	b.logger.Info("Channels created, signaling chain worker pool to start...")

	// Start pool of workers so they are ready for incoming messages from the kafka listeners
	b.startChainPool(ctx, wg)

	// Start pool kafka listeners
	b.startKafkaListenersPool(ctx, wg)
}

// getUsersByChain filters the users by the chain

// startChainPool starts the chain worker pool, creates required number of workers per chain,
// each worker will listen to a specific chain via a channel
func (b *Service) startChainPool(ctx context.Context, wg *sync.WaitGroup) {
	for _, worker := range b.workers {
		chain := worker.Chain
		kafkaMessageCh := b.chainChannelMap[chain]
		processor := b.processorMap[chain]

		for i := 0; i < worker.Total; i++ {
			workerNumber := i

			wg.Add(1)

			go b.chainWorker(ctx, workerNumber, processor, wg, kafkaMessageCh)
		}
	}
}

// Need to listen to the correct topics for each chain and forward it on the correct channel for picking
// up
func (b *Service) startKafkaListenersPool(ctx context.Context, wg *sync.WaitGroup) {
	for _, topic := range b.chainTopics {
		chain := topic.Chain
		kafkaMessageCh := b.chainChannelMap[chain]

		wg.Add(1)

		go b.kafkaListener(ctx, topic, wg, kafkaMessageCh)
	}
}

func (b *Service) chainWorker(
	ctx context.Context, workerID int, processor processors.Processor, wg *sync.WaitGroup,
	kafkaMessageCh chan kafka.Message,
) {
	for {
		select {
		case <-ctx.Done():
			b.logger.Info(fmt.Sprintf("Worker %d: Context Cancelled, stopping...", workerID))
			wg.Done()

			return
		case msg, ok := <-kafkaMessageCh:
			if !ok {
				b.logger.Info(fmt.Sprintf("Worker %d: No more tasks, exiting...", workerID))

				return
			}

			// b.logger.Info(fmt.Sprintf("Worker %d: Processing %s", workerID, msg))

			err := processor.Process(ctx, workerID, msg)
			if err != nil {
				b.logger.Error(fmt.Sprintf("Worker %d: error processing transactions, ", workerID), "error", err)
			}
		}
	}
}

// kafkaListener listens to the kafka topic and forwards the message to the channel for a specific chain
func (b *Service) kafkaListener(ctx context.Context, topic config.ChainTopicConfig, wg *sync.WaitGroup, kafkaMessageCh chan kafka.Message) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				b.logger.Info(fmt.Sprintf("Stopping kafka listener for `%s` chain, due to context cancellation...", topic.Chain))

				wg.Done()

				return
			default:
				err := b.kafkaAdapter.Subscribe(
					ctx, topic.Topic, func(topic string, key, value []byte) {
						// Process the message here
						b.logger.Debug("Received message", "topic", topic, "key", key)
						kafkaMessageCh <- kafka.Message{Topic: topic, Key: key, Value: value}
					},
					kafka.SubscriptionOptions{
						GroupID: topic.ConsumerGroup,
					},
				)
				if err != nil {
					b.logger.Info("Error received from subscription to kafka topic", err)
				}
			}
		}
	}()
}

// func (b *Service) processTransactions(ctx context.Context, workerID int, msg kafka.Message) error {
// 	var blockResponse response.Block
//
// 	err := json.Unmarshal(msg.Value, &blockResponse)
// 	if err != nil {
// 		return fmt.Errorf("worker %d: error marshalling block by number: %w", workerID, err)
// 	}
//
// 	for _, signature := range blockResponse.Signatures {
// 		transaction, err := b.bcProvider.GetTransaction(
// 			ctx, signature, &request.GetTransactionOptions{
// 				MaxSupportedTransactionVersion: 1,
// 			},
// 		)
// 		if err != nil {
// 			return fmt.Errorf("worker %d: error getting transaction: %w", workerID, err)
// 		}
//
// 		if transaction.Error != nil {
// 			return fmt.Errorf("worker %d: error returned from getting raw transaction: %v", workerID, transaction.Error)
// 		}
//
// 		value, err := json.Marshal(transaction.Result)
// 		if err != nil {
// 			return fmt.Errorf("worker %d: error marshalling raw transaction: %w", workerID, err)
// 		}
//
// 		// Here i will output to stdout
//
// 		// err = b.kafkaAdapter.PublishMessage(ctx, b.kafkaTransactionTopic, []byte(transaction.Result.Transaction.Signatures[0]), value)
// 		// if err != nil {
// 		// 	return fmt.Errorf("worker %d: error publishing message: %w", workerID, err)
// 		// }
//
// 		_ = value
// 	}
//
// 	return nil
// }
