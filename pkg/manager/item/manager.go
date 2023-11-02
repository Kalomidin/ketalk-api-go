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
	karatRepository     repository.KaratRepository
	categoryRepository  repository.CategoryRepository
	userPort            port.UserPort
	conversationPort    port.ConversationPort
	blobStorage         storage.AzureBlobStorage
}

func NewItemManager(itemRepository repository.ItemRepository, itemImageRepository repository.ItemImageRepository, userItemRepository repository.UserItemRepository, karatRepository repository.KaratRepository, categoryRepository repository.CategoryRepository, userPort port.UserPort, conversationPort port.ConversationPort, azureBlobStorage storage.AzureBlobStorage) ItemManager {
	return &itemManager{
		itemRepository,
		itemImageRepository,
		userItemRepository,
		karatRepository,
		categoryRepository,
		userPort,
		conversationPort,
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

	return m.repoItemIntoItemBlocks(ctx, items), nil
}

func (m *itemManager) GetItem(ctx context.Context, req GetItemRequest) (*GetItemResponse, error) {
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

	return &GetItemResponse{
		Item: Item{
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
		},
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
	} else {
		userItem.IsFavorite = req.IsFavorite
		if err := m.userItemRepository.Update(ctx, userItem); err != nil {
			return nil, err
		}
	}
	if req.IsFavorite {
		if err := m.itemRepository.IncrementFavoriteCount(ctx, req.ItemID); err != nil {
			return nil, err
		}
	} else {
		if err := m.itemRepository.DecrementFavoriteCount(ctx, req.ItemID); err != nil {
			return nil, err
		}
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
		if item.ItemStatus == string(ItemStatusSold) {
			// check if we have already a buyer
			_, err = m.userItemRepository.GetItemBuyer(ctx, req.ItemID)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			if err == nil {
				return nil, fmt.Errorf("item is already sold")
			}
		}
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

func (m *itemManager) IncrementConversationCount(ctx context.Context, req IncrementConversationCountRequest) error {
	return m.itemRepository.IncrementMessageCount(ctx, req.ItemID)
}

func (m *itemManager) GetAllKarats(ctx context.Context) ([]Karat, error) {
	karats, err := m.karatRepository.GetAllKarats(ctx)
	if err != nil {
		return nil, err
	}

	var resp []Karat = make([]Karat, len(karats))
	for i, karat := range karats {
		var locales map[string]KaratLocale = make(map[string]KaratLocale, len(karat.Locales))
		for locale, karatDescription := range karat.Locales {
			locales[locale] = KaratLocale{
				Name:        karatDescription.Name,
				Description: karatDescription.Description,
			}
		}
		resp[i] = Karat{
			ID:      karat.ID,
			Name:    karat.Name,
			Locales: locales,
		}
	}
	return resp, nil
}

func (m *itemManager) GetAllCategories(ctx context.Context) ([]Category, error) {
	categories, err := m.categoryRepository.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	var resp []Category = make([]Category, len(categories))
	for i, category := range categories {
		var locales map[string]CategoryLocale = make(map[string]CategoryLocale, len(category.Locales))
		for locale, categoryDescription := range category.Locales {
			locales[locale] = CategoryLocale{
				Name:        categoryDescription.Name,
				Description: categoryDescription.Description,
			}
		}
		resp[i] = Category{
			ID:      category.ID,
			Name:    category.Name,
			Locales: locales,
		}
	}
	return resp, nil
}

func (m *itemManager) GetSimilarItems(ctx context.Context, req GetSimilarItemsRequest) (*GetSimilarItemsResponse, error) {
	item, err := m.itemRepository.GetItem(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}

	var suggestedItems []ItemBlock = make([]ItemBlock, 0)
	var otherUserItems []ItemBlock = make([]ItemBlock, 0)
	// get suggested items if user is not owner of item
	if item.OwnerID != req.UserID {
		totalItemsCount := 20
		userOtherItemsCount := 2

		repoUserOtherItems, err := m.itemRepository.GetLimitedUserItems(ctx, item.OwnerID, userOtherItemsCount)
		if err != nil {
			return nil, err
		}
		otherUserItems = m.repoItemIntoItemBlocks(ctx, repoUserOtherItems)

		suggestedItemsCount := totalItemsCount - len(repoUserOtherItems)

		repoSuggestedItems, err := m.itemRepository.GetLimitedItemsByCategoryOrKarat(ctx, item.OwnerID, item.CategoryID, item.KaratID, suggestedItemsCount)
		if err != nil {
			return nil, err
		}
		suggestedItems = m.repoItemIntoItemBlocks(ctx, repoSuggestedItems)
	}

	return &GetSimilarItemsResponse{
		SuggestedItems: suggestedItems,
		OtherUserItems: otherUserItems,
	}, nil
}

func (m *itemManager) GetItemBuyers(ctx context.Context, req GetItemBuyersRequest) ([]ItemBuyer, error) {
	item, err := m.itemRepository.GetItem(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}

	conversations, err := m.conversationPort.GetItemConversations(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}
	var resp []ItemBuyer = make([]ItemBuyer, len(conversations))
	for i, conversation := range conversations {
		var member port.Member
		for _, m := range conversation.Members {
			if member.MemberID == item.OwnerID {
				continue
			}
			member = m
		}
		user, err := m.userPort.GetUser(ctx, member.MemberID)
		if err != nil {
			continue
		}
		var avatar *string
		if user.Image != nil {
			url := m.blobStorage.GetFrontDoorUrl(*user.Image)
			avatar = &url
		}
		resp[i] = ItemBuyer{
			ID:             conversation.ID,
			Name:           user.Username,
			Avatar:         avatar,
			LastMessagedAt: conversation.LastMessagedAt,
		}
	}
	return resp, nil
}

func (m *itemManager) CreatePurchase(ctx context.Context, req CreatePurchaseRequest) (*CreatePurchaseResponse, error) {
	userItem, err := m.userItemRepository.GetUserItem(ctx, req.BuyerID, req.ItemID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		userItem = &repository.UserItem{
			UserID:      req.BuyerID,
			ItemID:      req.ItemID,
			IsPurchased: true,
		}
		if err := m.userItemRepository.Insert(ctx, userItem); err != nil {
			return nil, err
		}
		return &CreatePurchaseResponse{
			ItemID:  req.ItemID,
			BuyerID: req.BuyerID,
		}, nil
	} else {
		userItem.IsPurchased = true
		if err := m.userItemRepository.Update(ctx, userItem); err != nil {
			return nil, err
		}
		return &CreatePurchaseResponse{
			ItemID:  req.ItemID,
			BuyerID: req.BuyerID,
		}, nil
	}
}

func (m *itemManager) repoItemIntoItemBlocks(ctx context.Context, repoItems []repository.Item) []ItemBlock {
	var userOtherItems []ItemBlock = make([]ItemBlock, len(repoItems))
	for i, item := range repoItems {
		image, err := m.itemImageRepository.GetItemThumbnail(ctx, item.ID)
		if err != nil {
			continue
		}
		thumbnail := m.blobStorage.GetFrontDoorUrl(image.Key)

		userOtherItems[i] = ItemBlock{
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
		}
	}
	return userOtherItems
}
