package port

import (
	"context"

	"github.com/google/uuid"
)

type GetItemRequest struct {
	ItemID uuid.UUID
}

type Item struct {
	ID      uuid.UUID
	Title   string
	Price   uint32
	OwnerID uuid.UUID
}

type ItemPort interface {
	GetItem(ctx context.Context, itemId uuid.UUID) (*Item, error)
	GetCovertImage(ctx context.Context, itemId uuid.UUID) (string, error)
	IncrementMessageCount(ctx context.Context, itemId uuid.UUID) error
}
