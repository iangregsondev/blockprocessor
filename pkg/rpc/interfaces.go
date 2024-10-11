package rpc

import "context"

type Client interface {
	Request(ctx context.Context, method string, params []interface{}) ([]byte, error)
}
