package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	iowrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/io"
	"github.com/iangregsondev/deblockprocessor/internal/wrappers/logger"
)

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPC struct {
	logger    logger.Logger
	ioWrapper iowrapper.IO

	rpcURL string
	apiKey string

	maxRetries          int
	retryDelay          time.Duration
	reportRetryAttempts bool
}

func NewRPC(logger logger.Logger, ioWrapper iowrapper.IO, rpcURL string, apiKey string, maxRetryOnError int, retryDelayMs int, reportRetryAttempts bool) *RPC {
	return &RPC{
		logger:              logger,
		ioWrapper:           ioWrapper,
		rpcURL:              rpcURL,
		apiKey:              apiKey,
		maxRetries:          maxRetryOnError,
		retryDelay:          time.Duration(retryDelayMs) * time.Millisecond,
		reportRetryAttempts: reportRetryAttempts,
	}
}

func (r *RPC) OldRequest(ctx context.Context, method string, params []interface{}) ([]byte, error) {
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

	body, err := r.ioWrapper.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
func (r *RPC) Request(ctx context.Context, method string, params []interface{}) ([]byte, error) {
	request := rpcRequest{
		Jsonrpc: "2.0",
		ID:      "go-client",
		Method:  method,
		Params:  params,
	}

	// Prepare the request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var result []byte

	currentRetries := 0

	// Use retry logic for the HTTP request
	err = retry.Do(
		func() error {
			// Create HTTP request
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.rpcURL, bytes.NewReader(requestBody))
			if err != nil {
				return fmt.Errorf("failed to create new HTTP request: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+r.apiKey)

			// Perform the HTTP request
			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("HTTP request failed: %w", err)
			}

			defer resp.Body.Close()

			// Check for non-OK status
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			// Read the response body
			body, err := r.ioWrapper.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			if currentRetries > 0 && r.reportRetryAttempts {
				r.logger.Warn("Retry successful, recovered", "totalAttempts", currentRetries)
			}

			result = body

			return nil
		},
		retry.Attempts(uint(r.maxRetries)),
		retry.Delay(r.retryDelay),
		retry.Context(ctx),
		retry.OnRetry(
			func(n uint, err error) {
				if r.reportRetryAttempts {
					r.logger.Warn("Retry", "attempt", n, "error", err)
				}

				currentRetries++
			},
		),
	)

	// Return error if all retry attempts fail
	if err != nil {
		return nil, fmt.Errorf("request failed after retries: %w", err)
	}

	// Return the successful response
	return result, nil
}
