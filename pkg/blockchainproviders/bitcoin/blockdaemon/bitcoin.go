package blockdaemon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/models/response"
	"github.com/iangregsondev/deblockprocessor/pkg/rpc"
)

type Provider struct {
	rpcClient rpc.Client
}

func NewProvider(rpcClient rpc.Client) *Provider {
	return &Provider{
		rpcClient: rpcClient,
	}
}

func (p *Provider) GetBlockHash(ctx context.Context, height int) (*response.GetBlockHashResponse, error) {
	result, err := p.rpcClient.Request(ctx, "getblockhash", []interface{}{height})
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash: %w", err)
	}

	var resp response.GetBlockHashResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block hash: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetBlockHeader(ctx context.Context, blockHash string) (*response.GetBlockHeaderResponse, error) {
	result, err := p.rpcClient.Request(ctx, "getblockheader", []interface{}{blockHash})
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash: %w", err)
	}

	var resp response.GetBlockHeaderResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block hash: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetBlockchainInfo(ctx context.Context) (*response.GetBlockchainInfoResponse, error) {
	result, err := p.rpcClient.Request(ctx, "getblockchaininfo", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash: %w", err)
	}

	var resp response.GetBlockchainInfoResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block hash: %w", err)
	}

	return &resp, nil
}

func (p *Provider) GetBlockCount(ctx context.Context) (*response.GetBlockCountResponse, error) {
	result, err := p.rpcClient.Request(ctx, "getblockcount", []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get block hash: %w", err)
	}

	var resp response.GetBlockCountResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse block hash: %w", err)
	}

	return &resp, nil
}
