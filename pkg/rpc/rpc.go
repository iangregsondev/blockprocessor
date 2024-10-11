package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPC struct {
	rpcURL string
	apiKey string
}

func NewRPC(rpcURL string, apiKey string) *RPC {
	return &RPC{
		rpcURL: rpcURL,
		apiKey: apiKey,
	}
}

func (r *RPC) Request(ctx context.Context, method string, params []interface{}) ([]byte, error) {
	request := rpcRequest{
		Jsonrpc: "2.0",
		ID:      "go-client",
		Method:  method,
		Params:  params,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.rpcURL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
