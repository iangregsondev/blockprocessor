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

func (p *Provider) GetBlockByHash(ctx context.Context, blockHash string, fullTransaction bool) (*response.BlockByHashResponse, error) {
	result, err := p.rpcClient.Request(ctx, "eth_getBlockByHash", []interface{}{blockHash, fullTransaction})
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %w", err)
	}

	var resp response.BlockByHashResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block by hash: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetTransactionByHash(ctx context.Context, transactionHash string) (*response.TransactionByHashResponse, error) {
	result, err := p.rpcClient.Request(ctx, "eth_getTransactionByHash", []interface{}{transactionHash})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}

	var resp response.TransactionByHashResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction by hash: %w", err)
	}

	return &resp, nil
}
