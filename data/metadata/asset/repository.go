package asset

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Repository struct {
	db           *gorm.DB
	alpacaClient *alpaca.Client
}

func NewRepository(db *gorm.DB, alpacaClient *alpaca.Client) (*Repository, error) {
	err := db.AutoMigrate(&alpaca.Asset{})

	if err != nil {
		return nil, err
	}

	var count int64

	result := db.Model(&alpaca.Asset{}).Count(&count)

	if result.Error != nil {
		return nil, result.Error
	}

	if count == 0 {
		log.Println("Populating assets")
		assets, err := alpacaClient.GetAssets(alpaca.GetAssetsRequest{
			Status: "active",
		})

		if err != nil {
			return nil, err
		}

		result = db.Model(&alpaca.Asset{}).CreateInBatches(assets, 100)

		if result.Error != nil {
			return nil, result.Error
		}

		log.Println("Populating assets...COMPLETED")
	}

	return &Repository{
		db:           db,
		alpacaClient: alpacaClient,
	}, nil
}

func (r *Repository) Get(symbol string) (*alpaca.Asset, error) {
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	var asset alpaca.Asset

	result := r.db.
		Model(&alpaca.Asset{}).
		Where("symbol = ?", symbol).
		First(&asset)

	if result.Error != nil {
		return nil, result.Error
	}

	return &asset, nil
}

func (r *Repository) GetAll() ([]alpaca.Asset, error) {
	var assets []alpaca.Asset
	result := r.db.Find(&assets)

	if result.Error != nil {
		return nil, result.Error
	}

	return assets, nil
}
