package item_manager

import (
	"context"
	"ketalk-api/pkg/manager/item/repository"
	"ketalk-api/pkg/manager/port"

	"github.com/google/uuid"
)

type itemPort struct {
	itemRepo      repository.ItemRepository
	itemImageRepo repository.ItemImageRepository
}

func NewItemPort(itemRepo repository.ItemRepository, itemImageRepo repository.ItemImageRepository) port.ItemPort {
	return &itemPort{
		itemRepo,
		itemImageRepo,
	}
}

func (p *itemPort) GetItem(ctx context.Context, itemID uuid.UUID) (*port.Item, error) {
	item, err := p.itemRepo.GetItem(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return &port.Item{
		ID:      item.ID,
		Title:   item.Title,
		Price:   item.Price,
		OwnerID: item.OwnerID,
	}, nil
}

func (p *itemPort) GetCovertImage(ctx context.Context, itemID uuid.UUID) (string, error) {
	image, err := p.itemImageRepo.GetItemThumbnail(ctx, itemID)
	if err != nil {
		return "", err
	}
	return image.Key, nil
}
