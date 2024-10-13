package bitcoin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
	"github.com/iangregsondev/deblockprocessor/internal/common/convert"
	"github.com/iangregsondev/deblockprocessor/internal/multichain/transactionobserver/models/config"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/models/response"
)

const (
	chain = "bitcoin"

	decimalPlaces = 8
)

type Processor struct {
	logger     logger.Logger
	bcProvider bitcoin.Provider

	users []config.UserConfig
}

func NewProcessor(logger logger.Logger, bcProvider bitcoin.Provider, users []config.UserConfig) *Processor {
	return &Processor{
		logger:     logger,
		bcProvider: bcProvider,
		users:      users,
	}
}

func (p *Processor) Process(ctx context.Context, workerID int, msg kafka.Message) error {
	var transaction response.RawTransaction

	err := json.Unmarshal(msg.Value, &transaction)
	if err != nil {
		return fmt.Errorf("error unmarshalling transaction: %w", err)
	}

	p.logger.Debug(fmt.Sprintf("[%s:worker %d] Processing transaction: %s", chain, workerID, transaction.Txid))

	// Prepare to track total input and output amounts for fee calculation
	var totalInputAmount, totalOutputAmount, userAmount float64

	var sourceAddresses, destinationAddresses []string

	involved := false
	involvedLocation := make([]string, 0)
	coinbase := false

	// Process the inputs
	for _, vin := range transaction.Vin {
		// handle coinbase transaction
		if vin.Coinbase != "" {
			coinbase = true

			// coinbase transactions have no traditional inputs
			// we still need to count these output values as the total amount being input (created) in the block reward.
			for _, vout := range transaction.Vout {
				totalInputAmount += vout.Value
			}

			// skip to the next iteration, coinbase should be the only input
			continue
		}

		// If it's not a coinbase transaction, fetch the previous transaction
		if vin.Txid != "" {
			// Fetch the previous transaction using the `txid` from Vin
			prevTransaction, err := p.bcProvider.GetRawTransaction(ctx, vin.Txid, true)
			if err != nil {
				return fmt.Errorf("error fetching previous transaction for txid %s: %w", vin.Txid, err)
			}

			// Use the `vout` index in Vin to get the correct output in the previous transaction
			if vin.Vout < len(prevTransaction.Result.Vout) {
				prevOutput := prevTransaction.Result.Vout[vin.Vout]
				sourceAddresses = append(sourceAddresses, prevOutput.ScriptPubKey.Address)
				totalInputAmount += prevOutput.Value

				// Check if the previous output address matches any user address
				for _, user := range p.users {
					if prevOutput.ScriptPubKey.Address == user.Address {
						involved = true

						involvedLocation = append(involvedLocation, "inputs")

						userAmount += prevOutput.Value
					}
				}
			} else {
				p.logger.Error(fmt.Sprintf("[%s:worker %d] Invalid vout index %d for previous transaction %s", chain, workerID, vin.Vout, vin.Txid))
			}
		}
	}

	// Process the outputs
	for _, vout := range transaction.Vout {
		destinationAddresses = append(destinationAddresses, vout.ScriptPubKey.Address)

		totalOutputAmount += vout.Value

		// Check if the address in this output matches any of the users
		for _, user := range p.users {
			if vout.ScriptPubKey.Address == user.Address {
				involved = true

				involvedLocation = append(involvedLocation, "outputs")

				userAmount += vout.Value
			}
		}
	}

	// Calculate the fees (Total Inputs - Total Outputs) for non-coinbase transactions
	fees := 0.0

	// coinbase transactions have no fees
	if totalInputAmount > 0 {
		fees = totalInputAmount - totalOutputAmount
	}

	if involved {
		kvPairs := []any{
			"txid", transaction.Txid,
			"involvedLocation", strings.Join(involvedLocation, ", "),
			"sourceAddresses", sourceAddresses,
			"destinationAddresses", destinationAddresses,
			"totalInputAmount", convert.FormatAmount(totalInputAmount, decimalPlaces),
			"totalOutputAmount", convert.FormatAmount(totalOutputAmount, decimalPlaces),
			"userAmount", convert.FormatAmount(userAmount, decimalPlaces),
			"fees", convert.FormatAmount(fees, decimalPlaces),
			"isCoinbase", coinbase,
		}

		p.logger.Info(fmt.Sprintf("[%s:worker %d] Transaction involving user detected:", chain, workerID), kvPairs...)
	}

	return nil
}
