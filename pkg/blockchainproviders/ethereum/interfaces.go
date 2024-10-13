package ethereum

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/models/response"
)

type Provider interface {
	GetBlockNumber(ctx context.Context) (*response.BlockNumberResponse, error)
	GetBlockByNumber(ctx context.Context, blockNumberHex string, fullTransaction bool) (*response.BlockByNumberResponse, error)
	GetBlockByHash(ctx context.Context, blockHash string, fullTransaction bool) (*response.BlockByHashResponse, error)
	GetTransactionByHash(ctx context.Context, transactionHash string) (*response.TransactionByHashResponse, error)
	GetTransactionReceipt(ctx context.Context, transactionHash string) (*response.TransactionReceiptResponse, error)
}
