package item_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetAllCategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID      uuid.UUID                 `json:"id"`
	Name    string                    `json:"name"`
	Locales map[string]CategoryLocale `json:"locales"`
}

type CategoryLocale struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *HttpHandler) GetAllCategories(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetAllCategories(ctx)
	return resp, err
}

func (h *handler) GetAllCategories(ctx *gin.Context) (*GetAllCategoriesResponse, error) {

	resp, err := h.manager.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	var categories []Category = make([]Category, len(resp))
	for i, category := range resp {
		var locales map[string]CategoryLocale = make(map[string]CategoryLocale, len(category.Locales))
		for j, locale := range category.Locales {
			locales[j] = CategoryLocale{
				Name:        locale.Name,
				Description: locale.Description,
			}
		}
		categories[i] = Category{
			ID:      category.ID,
			Name:    category.Name,
			Locales: locales,
		}
	}
	return &GetAllCategoriesResponse{
		Categories: categories,
	}, nil
}
