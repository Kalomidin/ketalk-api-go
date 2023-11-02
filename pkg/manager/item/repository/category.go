package repository

import (
	"context"

	"gorm.io/gorm"
)

type categoryRepository struct {
	*gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db,
	}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]Category, error) {
	var categories []Category = make([]Category, 0)
	resp := r.Find(&categories)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return categories, nil
}

func (r *categoryRepository) Migrate() error {
	if err := r.AutoMigrate(&Category{}); err != nil {
		return err
	}
	categories, err := r.GetAllCategories(context.Background())
	if err != nil {
		return err
	}
	if len(categories) > 0 {
		return nil
	}

	return r.CreateInBatches([]Category{
		{
			Locales: map[string]CategoryLocale{
				"en": {
					Description: "Worn around the neck.",
					Name:        "Rings",
				},
			},
			Name: "Rings",
		},
		{
			Locales: map[string]CategoryLocale{
				"en": {
					Description: "Worn around the wrist.",
					Name:        "Bracelets",
				},
			},
			Name: "Bracelets",
		},
		{
			Locales: map[string]CategoryLocale{
				"en": {
					Description: "Pendant with photo space.",
					Name:        "Lockets",
				},
			},
			Name: "Lockets",
		},
		{
			Locales: map[string]CategoryLocale{
				"en": {
					Description: "Decorative chain for ankle.",
					Name:        "Anklets",
				},
			},
			Name: "Anklets",
		},
	}, 4).Error
}
