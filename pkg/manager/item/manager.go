package item_manager

import (
	"context"
	"fmt"
	"ketalk-api/pkg/manager/item/repository"
	"ketalk-api/storage"
	"time"
)

type itemManager struct {
	itemRepository      repository.ItemRepository
	itemImageRepository repository.ItemImageRepository
	blobStorage         storage.AzureBlobStorage
}

func NewItemManager(itemRepository repository.ItemRepository, itemImageRepository repository.ItemImageRepository, azureBlobStorage storage.AzureBlobStorage) ItemManager {
	return &itemManager{
		itemRepository,
		itemImageRepository,
		azureBlobStorage,
	}
}

func (m *itemManager) AddItem(ctx context.Context, item AddItemRequest) (*AddItemResponse, error) {
	repoItem := repository.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Size:        item.Size,
		Weight:      item.Weight,
		OwnerID:     item.OwnerID,
		Negotiable:  item.Negotiable,
		KaratID:     item.KaratID,
		CategoryID:  item.CategoryID,
		GeofenceID:  item.GeofenceID,
		ItemStatus:  string(ItemStatusActive),
	}
	err := m.itemRepository.AddItem(ctx, &repoItem)
	if err != nil {
		return nil, err
	}
	var images []repository.ItemImage = make([]repository.ItemImage, len(item.Images))
	for i, image := range item.Images {
		key := fmt.Sprintf("%+v_%s", time.Now().UTC().UnixNano(), image)
		images[i] = repository.ItemImage{
			Key:     key,
			ItemID:  repoItem.ID,
			IsCover: image == item.Thumbnail,
		}
	}
	if err = m.itemImageRepository.AddItemImages(ctx, repoItem.ID, images); err != nil {
		return nil, err
	}

	var generatedUrls []SignedUrlWithImageID = make([]SignedUrlWithImageID, 0)
	for i, image := range images {
		url, err := m.blobStorage.GeneratePresignedUrlToUpload(image.Key)
		if err != nil {
			continue
		}
		generatedUrls = append(generatedUrls, SignedUrlWithImageID{
			ID:        image.ID,
			SignedUrl: url,
			Name:      item.Images[i],
		})
	}

	return &AddItemResponse{
		ID:            repoItem.ID,
		CreatedAt:     repoItem.CreatedAt,
		PresignedUrls: generatedUrls,
	}, nil
}

func (m *itemManager) GetItems(ctx context.Context, req GetItemsRequest) ([]Item, error) {
	items, err := m.itemRepository.GetItems(ctx, req.GeofenceID)
	if err != nil {
		return nil, err
	}
	var resp []Item
	for _, item := range items {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail, err := m.blobStorage.GeneratePresignedUrlToRead(image.Key)
		if err != nil {
			continue
		}
		resp = append(resp, Item{
			ID:            item.ID,
			Title:         item.Title,
			Description:   item.Description,
			Price:         item.Price,
			OwnerID:       item.OwnerID,
			FavoriteCount: item.FavoriteCount,
			MessageCount:  item.MessageCount,
			SeenCount:     item.SeenCount,
			ItemStatus:    ItemStatus(item.ItemStatus),
			CreatedAt:     item.CreatedAt,
			Thumbnail:     thumbnail,
		})
	}
	return resp, nil
}

func (m *itemManager) UploadItemImages(ctx context.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error) {
	if err := m.itemImageRepository.UpdateItemImagesToUploaded(ctx, r.ItemID, r.ImageIds); err != nil {
		return nil, err
	}

	return &UploadItemImagesResponse{}, nil
}
