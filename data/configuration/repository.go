package configuration

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

const key = "justme"

func NewRepository(db *gorm.DB) (*Repository, error) {
	err := db.AutoMigrate(&AppConfiguration{})

	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) IsEmpty() (bool, error) {
	var count int64
	result := r.db.Model(&AppConfiguration{}).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count == 0, nil
}

func (r *Repository) Get() (*AppConfiguration, error) {
	config := AppConfiguration{
		Key:    key,
		X:      0,
		Y:      0,
		Width:  1024,
		Height: 768,
	}

	result := r.db.FirstOrCreate(&config)

	if result.Error != nil {
		return nil, result.Error
	}

	return &config, nil
}

func (r *Repository) UpdateWindow(x, y, width, height int) error {
	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "key"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"x":      x,
			"y":      y,
			"width":  width,
			"height": height,
		}),
	}).Create(&AppConfiguration{
		Key:    key,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	})

	return result.Error
}
