package blockdaemon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/rpc"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana/models/request"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana/models/response"
)

type Provider struct {
	rpcClient rpc.Client
}

func NewProvider(rpcClient rpc.Client) *Provider {
	return &Provider{
		rpcClient: rpcClient,
	}
}

func (p *Provider) GetBlockHeight(ctx context.Context) (*response.BlockHeightResponse, error) {
	result, err := p.rpcClient.Request(ctx, "getBlockHeight", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get block height: %w", err)
	}

	var resp response.BlockHeightResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block height: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetBlock(ctx context.Context, blockNumber int64, options *request.GetBlockOptions) (*response.BlockResponse, error) {
	// Set default values
	defaultOptions := request.GetBlockOptions{
		Encoding:                       "json",
		TransactionDetails:             "full",
		Rewards:                        false,
		MaxSupportedTransactionVersion: 0,
	}

	// If options are provided, overwrite default values with what is passed in
	if options != nil {
		if options.Encoding != "" {
			defaultOptions.Encoding = options.Encoding
		}

		if options.TransactionDetails != "" {
			defaultOptions.TransactionDetails = options.TransactionDetails
		}

		defaultOptions.Rewards = options.Rewards

		if options.MaxSupportedTransactionVersion != 0 {
			defaultOptions.MaxSupportedTransactionVersion = options.MaxSupportedTransactionVersion
		}
	}

	params := []interface{}{
		blockNumber,
		map[string]interface{}{
			"encoding":                       defaultOptions.Encoding,
			"transactionDetails":             defaultOptions.TransactionDetails,
			"rewards":                        defaultOptions.Rewards,
			"maxSupportedTransactionVersion": defaultOptions.MaxSupportedTransactionVersion,
		},
	}

	result, err := p.rpcClient.Request(ctx, "getBlock", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockNumberHex by number: %w", err)
	}

	var resp response.BlockResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse blockNumberHex by number: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetTransaction(
	ctx context.Context, transactionSignature string, options *request.GetTransactionOptions,
) (*response.TransactionResponse, error) {
	// Set default values
	defaultOptions := request.GetTransactionOptions{
		Commitment:                     "finalized",
		Encoding:                       "json",
		MaxSupportedTransactionVersion: 1,
	}

	// If options are provided, overwrite default values with what is passed in
	if options != nil {
		if options.Commitment != "" {
			defaultOptions.Commitment = options.Commitment
		}

		if options.Encoding != "" {
			defaultOptions.Encoding = options.Encoding
		}

		if options.MaxSupportedTransactionVersion != 0 {
			defaultOptions.MaxSupportedTransactionVersion = options.MaxSupportedTransactionVersion
		}
	}

	// Prepare the request params for Solana's getTransaction method
	params := []interface{}{
		transactionSignature,
		map[string]interface{}{
			"commitment":                     defaultOptions.Commitment,
			"encoding":                       defaultOptions.Encoding,
			"maxSupportedTransactionVersion": defaultOptions.MaxSupportedTransactionVersion,
		},
	}

	result, err := p.rpcClient.Request(ctx, "getTransaction", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var resp response.TransactionResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction: %w", err)
	}

	return &resp, nil
}
