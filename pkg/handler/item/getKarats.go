package item_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetAllKaratsResponse struct {
	Karats []Karat `json:"karats"`
}

type Karat struct {
	ID      uuid.UUID              `json:"id"`
	Name    string                 `json:"name"`
	Locales map[string]KaratLocale `json:"locales"`
}

type KaratLocale struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *HttpHandler) GetAllKarats(ctx *gin.Context, r *http.Request) (interface{}, error) {
	resp, err := h.handler.GetAllKarats(ctx)
	return resp, err
}

func (h *handler) GetAllKarats(ctx *gin.Context) (*GetAllKaratsResponse, error) {

	resp, err := h.manager.GetAllKarats(ctx)
	if err != nil {
		return nil, err
	}

	var karats []Karat = make([]Karat, len(resp))
	for i, karat := range resp {
		var locales map[string]KaratLocale = make(map[string]KaratLocale, len(karat.Locales))
		for j, locale := range karat.Locales {
			locales[j] = KaratLocale{
				Name:        locale.Name,
				Description: locale.Description,
			}
		}
		karats[i] = Karat{
			ID:      karat.ID,
			Name:    karat.Name,
			Locales: locales,
		}
	}
	return &GetAllKaratsResponse{
		Karats: karats,
	}, nil
}
