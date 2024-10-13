package blockdaemon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/rpc"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/models/response"
)

type Provider struct {
	rpcClient rpc.Client
}

func NewProvider(rpcClient rpc.Client) *Provider {
	return &Provider{
		rpcClient: rpcClient,
	}
}

func (p *Provider) GetBlockNumber(ctx context.Context) (*response.BlockNumberResponse, error) {
	result, err := p.rpcClient.Request(ctx, "eth_blockNumber", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get block number: %w", err)
	}

	var resp response.BlockNumberResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block number: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetBlockByNumber(ctx context.Context, blockNumberHex string, fullTransaction bool) (*response.BlockByNumberResponse, error) {
	result, err := p.rpcClient.Request(ctx, "eth_getBlockByNumber", []interface{}{blockNumberHex, fullTransaction})
	if err != nil {
		return nil, fmt.Errorf("failed to get blockNumberHex by number: %w", err)
	}

	var resp response.BlockByNumberResponse
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

func (p *Provider) GetTransactionReceipt(ctx context.Context, transactionHash string) (*response.TransactionReceiptResponse, error) {
	result, err := p.rpcClient.Request(ctx, "eth_getTransactionReceipt", []interface{}{transactionHash})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}

	var resp response.TransactionReceiptResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction by hash: %w", err)
	}

	return &resp, nil
}
