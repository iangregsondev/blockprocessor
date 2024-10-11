package signalutils

import (
	"context"
	"os"
	"os/signal"

	"github.com/iangregsondev/deblockprocessor/pkg/logger"
)

type Worker func(ctx context.Context)

func BuildSignalHandler(logger logger.Logger, cancelFunc context.CancelFunc, sigTypes ...os.Signal) Worker {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sigTypes...)

	return func(ctx context.Context) {
		go func() {
			select {
			case <-ctx.Done():
				return
			case sig := <-sigs:
				_ = sig
				logger.Warn(
					"system signal received: starting graceful shutdown",
					"signal", sig,
				)
				// Call the cancel function to trigger a graceful shutdown
				cancelFunc()
			}
		}()
	}
}
