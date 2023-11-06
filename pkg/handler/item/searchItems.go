package item_handler

import (
	"fmt"
	item_manager "ketalk-api/pkg/manager/item"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SearchItemsResponse struct {
	Items []ItemBlock `json:"items"`
}

// search?priceRange=100,1000&karatIds=18,24&categoryIds=ring,necklace&sizeRange=10,20&keyword=hello

func (h *HttpHandler) SearchItems(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.SearchItems(ctx)
	return resp, err
}

func (h *handler) SearchItems(ctx *gin.Context) (*SearchItemsResponse, error) {
	priceRange, err := getRangeUint(ctx, "priceRange")
	if err != nil {
		return nil, err
	}
	sizeRange, err := getRangeFloat32(ctx, "sizeRange")
	if err != nil {
		return nil, err
	}
	karatIds, err := getIds(ctx, "karatIds")
	if err != nil {
		return nil, err
	}
	categoryIds, err := getIds(ctx, "categoryIds")
	if err != nil {
		return nil, err
	}
	keyword := ctx.Query("keyword")
	manReq := item_manager.SearchItemsRequest{
		Keyword:     keyword,
		PriceRange:  priceRange,
		SizeRange:   sizeRange,
		KaratIDs:    karatIds,
		CategoryIDs: categoryIds,
	}
	itemBlocks, err := h.manager.SearchItems(ctx, manReq)
	if err != nil {
		return nil, err
	}
	var items []ItemBlock = make([]ItemBlock, len(itemBlocks))
	for i, item := range itemBlocks {
		items[i] = ItemBlock{
			ID:            item.ID,
			Title:         item.Title,
			Description:   item.Description,
			Price:         item.Price,
			OwnerID:       item.OwnerID,
			FavoriteCount: item.FavoriteCount,
			MessageCount:  item.MessageCount,
			SeenCount:     item.SeenCount,
			ItemStatus:    string(item.ItemStatus),
			CreatedAt:     item.CreatedAt.UTC().Unix(),
			Thumbnail:     item.Thumbnail,
			IsHidden:      item.IsHidden,
		}
	}

	return &SearchItemsResponse{
		Items: items,
	}, nil
}

func getRangeUint(ctx *gin.Context, key string) ([]uint32, error) {
	rangeString := ctx.Query(key)
	if rangeString == "" {
		return []uint32{0, math.MaxUint32}, nil
	}

	rng := strings.Split(rangeString, ",")
	fmt.Printf("rng: %v, len: %+v, first: %+v, rangeString: %+v\n", rng, len(rng), rng[0], rangeString)
	if len(rng) == 0 {
		return []uint32{0, math.MaxUint32}, nil
	} else if len(rng) == 1 {
		minVal, ok := strconv.ParseUint(rng[0], 10, 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		return []uint32{uint32(minVal), math.MaxUint32}, nil
	} else if len(rng) == 2 {
		minVal, ok := strconv.ParseUint(rng[0], 10, 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for first value for: %s", key)
		}
		maxVal, ok := strconv.ParseUint(rng[1], 10, 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		if minVal > maxVal {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		return []uint32{uint32(minVal), uint32(maxVal)}, nil
	} else {
		return []uint32{0, math.MaxUint32}, nil
	}
}

func getRangeFloat32(ctx *gin.Context, key string) ([]float32, error) {
	rangeString := ctx.Query(key)
	if rangeString == "" {
		return []float32{0, math.MaxFloat32}, nil
	}
	rng := strings.Split(rangeString, ",")
	if len(rng) == 1 {
		minVal, ok := strconv.ParseFloat(rng[0], 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		return []float32{float32(minVal), math.MaxFloat32}, nil
	} else if len(rng) == 2 {
		minVal, ok := strconv.ParseFloat(rng[0], 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		maxVal, ok := strconv.ParseFloat(rng[1], 32)
		if ok != nil {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		if minVal > maxVal {
			return nil, fmt.Errorf("invalid param or key for: %s", key)
		}
		return []float32{float32(minVal), float32(maxVal)}, nil
	} else {
		return []float32{0, math.MaxFloat32}, nil
	}
}

func getIds(ctx *gin.Context, key string) ([]uuid.UUID, error) {
	idsString := ctx.Query(key)
	if idsString == "" {
		return []uuid.UUID{}, nil
	}
	ids := strings.Split(idsString, ",")
	fmt.Printf("ids: %v, len: %+v, first: %+v, idsString: %+v\n", ids, len(ids), ids[0], idsString)
	var uuids []uuid.UUID
	for _, id := range ids {
		_uuid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, _uuid)
	}
	return uuids, nil
}
