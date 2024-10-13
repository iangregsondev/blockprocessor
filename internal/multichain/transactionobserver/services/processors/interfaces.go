package processors

import (
	"context"

	"github.com/iangregsondev/deblockprocessor/internal/adapters/kafka"
)

type Processor interface {
	Process(ctx context.Context, workerID int, msg kafka.Message) error
}
