package item_manager

import (
	"context"
	"fmt"
	"ketalk-api/pkg/manager/item/repository"
	"ketalk-api/pkg/manager/port"
	"ketalk-api/storage"
	"time"

	"gorm.io/gorm"
)

type itemManager struct {
	itemRepository      repository.ItemRepository
	itemImageRepository repository.ItemImageRepository
	userItemRepository  repository.UserItemRepository
	userPort            port.UserPort
	blobStorage         storage.AzureBlobStorage
}

func NewItemManager(itemRepository repository.ItemRepository, itemImageRepository repository.ItemImageRepository, userItemRepository repository.UserItemRepository, userPort port.UserPort, azureBlobStorage storage.AzureBlobStorage) ItemManager {
	return &itemManager{
		itemRepository,
		itemImageRepository,
		userItemRepository,
		userPort,
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

func (m *itemManager) GetItems(ctx context.Context, req GetItemsRequest) ([]ItemBlock, error) {
	items, err := m.itemRepository.GetItems(ctx, req.GeofenceID, req.UserID)
	if err != nil {
		return nil, err
	}
	var resp []ItemBlock
	for _, item := range items {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail := m.blobStorage.GetFrontDoorUrl(image.Key)

		resp = append(resp, ItemBlock{
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
			IsHidden:      item.IsHidden,
		})
	}
	return resp, nil
}

func (m *itemManager) GetItem(ctx context.Context, req GetItemRequest) (*Item, error) {
	item, err := m.itemRepository.GetItem(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}
	if item.IsHidden && item.OwnerID != req.UserID {
		return nil, fmt.Errorf("item is hidden")
	}

	itemImages, err := m.itemImageRepository.GetItemImages(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	var images []string = make([]string, len(itemImages))
	var thumbnail string
	for i, image := range itemImages {
		url := m.blobStorage.GetFrontDoorUrl(image.Key)
		if image.IsCover {
			thumbnail = url
		}
		images[i] = url
	}

	owner, err := m.userPort.GetUser(ctx, item.OwnerID)
	if err != nil {
		return nil, err
	}
	var ownerAvatar *string
	if owner.Image != nil {
		url := m.blobStorage.GetFrontDoorUrl(*owner.Image)
		ownerAvatar = &url
	}
	isUserFavorite := false
	userItem, err := m.userItemRepository.GetUserItem(ctx, req.UserID, item.ID)
	if err == nil && userItem != nil {
		isUserFavorite = userItem.IsFavorite
	}
	return &Item{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Owner: ItemOwner{
			ID:     owner.ID,
			Name:   owner.Username,
			Avatar: ownerAvatar,
		},
		FavoriteCount:  item.FavoriteCount,
		MessageCount:   item.MessageCount,
		SeenCount:      item.SeenCount,
		ItemStatus:     ItemStatus(item.ItemStatus),
		CreatedAt:      item.CreatedAt,
		Thumbnail:      thumbnail,
		Images:         images,
		IsUserFavorite: isUserFavorite,
		IsHidden:       item.IsHidden,
		Negotiable:     item.Negotiable,
	}, nil
}

func (m *itemManager) UploadItemImages(ctx context.Context, r UploadItemImagesRequest) (*UploadItemImagesResponse, error) {
	if err := m.itemImageRepository.UpdateItemImagesToUploaded(ctx, r.ItemID, r.ImageIds); err != nil {
		return nil, err
	}

	return &UploadItemImagesResponse{}, nil
}

func (m *itemManager) GetFavoriteItems(ctx context.Context, r GetFavoriteItemsRequest) ([]ItemBlock, error) {
	items, err := m.userItemRepository.GetUserFavoriteItems(ctx, r.UserID)
	if err != nil {
		return nil, err
	}
	var resp []ItemBlock
	for _, item := range items {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail := m.blobStorage.GetFrontDoorUrl(image.Key)

		resp = append(resp, ItemBlock{
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
			IsHidden:      item.IsHidden,
		})
	}
	return resp, nil
}

func (m *itemManager) FavoriteItem(ctx context.Context, req FavoriteItemRequest) (*FavoriteItemResponse, error) {
	userItem, err := m.userItemRepository.GetUserItem(ctx, req.UserID, req.ItemID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		userItem = &repository.UserItem{
			UserID:     req.UserID,
			ItemID:     req.ItemID,
			IsFavorite: req.IsFavorite,
		}
		if err := m.userItemRepository.Insert(ctx, userItem); err != nil {
			return nil, err
		}
		return &FavoriteItemResponse{}, nil
	}
	userItem.IsFavorite = req.IsFavorite
	if err := m.userItemRepository.Update(ctx, userItem); err != nil {
		return nil, err
	}
	return &FavoriteItemResponse{}, nil
}

func (m *itemManager) GetUserItems(ctx context.Context, req GetUserItemsRequest) ([]ItemBlock, error) {
	items, err := m.itemRepository.GetUserItems(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	var resp []ItemBlock
	for _, item := range items {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail := m.blobStorage.GetFrontDoorUrl(image.Key)

		resp = append(resp, ItemBlock{
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
			IsHidden:      item.IsHidden,
		})
	}
	return resp, nil
}

func (m *itemManager) GetPurchasedItems(ctx context.Context, req GetPurchasedItemsRequest) ([]ItemBlock, error) {
	items, err := m.userItemRepository.GetPurchasedItems(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	var resp []ItemBlock
	for _, item := range items {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail := m.blobStorage.GetFrontDoorUrl(image.Key)

		resp = append(resp, ItemBlock{
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
			IsHidden:      item.IsHidden,
		})
	}
	return resp, nil
}

func (m *itemManager) UpdateItem(ctx context.Context, req UpdateItemRequest) (*UpdateItemResponse, error) {
	item, err := m.itemRepository.GetItem(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}

	if item.OwnerID != req.UserID {
		return nil, fmt.Errorf("user is not owner of item")
	}
	if req.IsHidden != nil {
		item.IsHidden = *req.IsHidden
	}
	if req.ItemStatus != nil {
		item.ItemStatus = string(*req.ItemStatus)
	}
	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Negotiable != nil {
		item.Negotiable = *req.Negotiable
	}
	if err := m.itemRepository.Update(ctx, item); err != nil {
		return nil, err
	}
	return &UpdateItemResponse{}, nil
}
