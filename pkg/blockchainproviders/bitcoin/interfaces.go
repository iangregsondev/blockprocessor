package bitcoin

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/models/response"
)

type Provider interface {
	GetBlockCount(ctx context.Context) (*response.GetBlockCountResponse, error)
	GetBlockchainInfo(ctx context.Context) (*response.GetBlockchainInfoResponse, error)
	GetBlockHash(ctx context.Context, height int) (*response.GetBlockHashResponse, error)
	GetBlock(ctx context.Context, blockHash string) (*response.GetBlockResponse, error)
	GetBlockHeader(ctx context.Context, blockHash string) (*response.GetBlockHeaderResponse, error)
	GetRawTransaction(ctx context.Context, txID string, verbose bool) (*response.GetRawTransactionResponse, error)
}
