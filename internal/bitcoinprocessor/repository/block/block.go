package block

import (
	"errors"

	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/adapters/sqlite"
	"github.com/iangregsondev/deblockprocessor/internal/bitcoinprocessor/entities"
	"gorm.io/gorm"
)

var ErrBlockHeightNotFound = errors.New("block height not found")

type Repository struct {
	dbAdapter sqlite.Database
}

func NewBlockRepository(dbAdapter sqlite.Database) *Repository {
	return &Repository{dbAdapter: dbAdapter}
}

// CreateOrUpdateBlockHeight creates a new BlockHeight if it doesn't exist,
// or updates the existing record.
func (r *Repository) CreateOrUpdateBlockHeight(height int) error {
	db := r.dbAdapter.GetDB()

	var blockHeight entities.BlockHeight
	err := db.First(&blockHeight).Error

	if err == nil {
		// Update existing record
		blockHeight.Height = height

		return db.Save(&blockHeight).Error
	}

	// Create new record
	blockHeight = entities.BlockHeight{Height: height}

	return db.Create(&blockHeight).Error
}

// GetBlockHeight retrieves the current BlockHeight record if it exists.
func (r *Repository) GetBlockHeight() (*entities.BlockHeight, error) {
	db := r.dbAdapter.GetDB()

	var blockHeight entities.BlockHeight

	err := db.First(&blockHeight).Error
	if err != nil {
		// Check if the error is "record not found", in which case return nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBlockHeightNotFound
		}

		return nil, err
	}

	return &blockHeight, nil
}

// DeleteBlockHeight deletes the existing BlockHeight record if it exists.
func (r *Repository) DeleteBlockHeight() error {
	db := r.dbAdapter.GetDB()

	return db.Delete(&entities.BlockHeight{}).Error
}
