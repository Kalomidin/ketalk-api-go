package repository

import (
	"context"

	"gorm.io/gorm"
)

type karatRepository struct {
	*gorm.DB
}

func NewKaratRepository(db *gorm.DB) KaratRepository {
	return &karatRepository{
		db,
	}
}

func (r *karatRepository) GetAllKarats(ctx context.Context) ([]Karat, error) {
	var karats []Karat = make([]Karat, 0)
	resp := r.Find(&karats)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return karats, nil
}

func (r *karatRepository) Migrate() error {
	if err := r.AutoMigrate(&Karat{}); err != nil {
		return err
	}
	karats, err := r.GetAllKarats(context.Background())
	if err != nil {
		return err
	}
	if len(karats) > 0 {
		return nil
	}
	return r.CreateInBatches([]Karat{
		{
			Locales: map[string]KaratLocale{
				"en": {
					Description: "58.3% pure gold, alloyed mix.",
					Name:        "14K",
				},
			},
			Name: "14K",
		},
		{
			Locales: map[string]KaratLocale{
				"en": {
					Description: "75% pure gold, stronger shine.",
					Name:        "18K",
				},
			},
			Name: "18K",
		},
		{
			Locales: map[string]KaratLocale{
				"en": {
					Description: "91.7% pure gold, softer texture.",
					Name:        "22K",
				},
			},
			Name: "22K",
		},
		{
			Locales: map[string]KaratLocale{
				"en": {
					Description: "100% pure gold, soft & shiny.",
					Name:        "24K",
				},
			},
			Name: "24K",
		},
	}, 4).Error
}
