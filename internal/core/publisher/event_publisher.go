package publisher

import (
	"context"

	"github.com/salobato/ordermanager/internal/core/entity"
)

type EventPublisher interface {
	PublishOrderStatusChanged(ctx context.Context, event entity.OrderEvent) error
}
