package bitcoin

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/bitcoin/models/response"
)

type Provider interface {
	GetBlockCount(ctx context.Context) (*response.GetBlockCountResponse, error)
	GetBlockchainInfo(ctx context.Context) (*response.GetBlockchainInfoResponse, error)
	GetBlockHash(ctx context.Context, height int) (*response.GetBlockHashResponse, error)
	GetBlockHeader(ctx context.Context, blockHash string) (*response.GetBlockHeaderResponse, error)
}
