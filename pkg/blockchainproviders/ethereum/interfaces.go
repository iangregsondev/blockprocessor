package ethereum

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/pkg/blockchainproviders/ethereum/models/response"
)

type Provider interface {
	GetBlockNumber(ctx context.Context) (*response.BlockNumberResponse, error)
	GetBlockByNumber(ctx context.Context, blockNumberHex string, fullTransaction bool) (*response.BlockByNumberResponse, error)
}
