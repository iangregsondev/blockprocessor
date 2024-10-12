package solana

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana/models/request"
	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/solana/models/response"
)

type Provider interface {
	GetBlockHeight(ctx context.Context) (*response.BlockHeightResponse, error)
	GetBlock(ctx context.Context, blockNumber int64, options *request.GetBlockOptions) (*response.BlockResponse, error)

	GetBlockByHash(ctx context.Context, blockHash string, fullTransaction bool) (*response.BlockByHashResponse, error)
	GetTransactionByHash(ctx context.Context, transactionHash string) (*response.TransactionByHashResponse, error)
}
