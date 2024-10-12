package block

import (
	"errors"

	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/adapters/sqlite"
	"github.com/iangregsondev/deblockprocessor/internal/solana/blockprocessor/entities"
	"gorm.io/gorm"
)

var ErrBlockNumberNotFound = errors.New("block number not found")

type Repository struct {
	dbAdapter sqlite.Database
}

func NewBlockRepository(dbAdapter sqlite.Database) *Repository {
	return &Repository{dbAdapter: dbAdapter}
}

// CreateOrUpdateBlockNumber creates a new BlockNumber if it doesn't exist,
// or updates the existing record.
func (r *Repository) CreateOrUpdateBlockNumber(height int64) error {
	db := r.dbAdapter.GetDB()

	var blockHeight entities.BlockHeight
	err := db.First(&blockHeight).Error

	if err == nil {
		// Update existing record
		blockHeight.BlockHeight = height

		return db.Save(&blockHeight).Error
	}

	// Create new record
	blockHeight = entities.BlockHeight{BlockHeight: height}

	return db.Create(&blockHeight).Error
}

// GetBlockHeight retrieves the current BlockNumber record if it exists.
func (r *Repository) GetBlockHeight() (*entities.BlockHeight, error) {
	db := r.dbAdapter.GetDB()

	var blockHeight entities.BlockHeight

	err := db.First(&blockHeight).Error
	if err != nil {
		// Check if the error is "record not found", in which case return nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBlockNumberNotFound
		}

		return nil, err
	}

	return &blockHeight, nil
}

// DeleteBlockHeight deletes the existing BlockNumber record if it exists.
func (r *Repository) DeleteBlockHeight() error {
	db := r.dbAdapter.GetDB()

	return db.Delete(&entities.BlockHeight{}).Error
}
