package ethereum

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/common/convert"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/models/response"
)

const (
	chain = "ethereum"

	decimalPlaces = 18
)

type Processor struct {
	logger     logger.Logger
	bcProvider ethereum.Provider

	users []config.UserConfig
}

func NewProcessor(logger logger.Logger, bcProvider ethereum.Provider, users []config.UserConfig) *Processor {
	return &Processor{
		logger:     logger,
		bcProvider: bcProvider,
		users:      users,
	}
}

func (p *Processor) Process(ctx context.Context, workerID int, msg kafka.Message) error {
	var transaction response.Transaction

	err := json.Unmarshal(msg.Value, &transaction)
	if err != nil {
		return err
	}

	involved := false
	involvedLocation := make([]string, 0)

	for _, user := range p.users {
		if strings.EqualFold(transaction.From, user.Address) {
			involved = true

			involvedLocation = append(involvedLocation, "From")

			break
		}

		if strings.EqualFold(transaction.To, user.Address) {
			involved = true

			involvedLocation = append(involvedLocation, "To")

			break
		}
	}

	if involved {
		value, err := convert.HexToBigInt(transaction.Value)
		if err != nil {
			return fmt.Errorf("error converting value to decimal: %w", err)
		}

		ethValue := convert.WeiToEthUsingBigInt(value)

		// Need to get the receipt to get the actual gas that was used
		receipt, err := p.bcProvider.GetTransactionReceipt(ctx, transaction.Hash)
		if err != nil {
			return fmt.Errorf("error getting transaction receipt: %w", err)
		}

		gasUsed, err := convert.HexToBigInt(receipt.Result.GasUsed)
		if err != nil {
			return fmt.Errorf("error converting gas used to decimal: %w", err)
		}

		gasPrice, err := convert.HexToBigInt(transaction.GasPrice)
		if err != nil {
			return fmt.Errorf("error converting gas price to decimal: %w", err)
		}

		// Calculate fees
		gasTotal := new(big.Int).Mul(gasUsed, gasPrice)

		ethGasTotal := convert.WeiToEthUsingBigInt(gasTotal)

		kvPairs := []any{
			"transactionHash", transaction.Hash,
			"involvedLocation", strings.Join(involvedLocation, ", "),
			"fromAddress", transaction.From,
			"toAddress", transaction.To,
			"value", convert.FormatAmountBigFloat(ethValue, decimalPlaces),
			"gasPrice", convert.FormatAmountBigFloat(ethGasTotal, decimalPlaces),
		}

		p.logger.Info(fmt.Sprintf("[%s:worker %d] Transaction involving user detected:", chain, workerID), kvPairs...)
	}

	return nil
}
