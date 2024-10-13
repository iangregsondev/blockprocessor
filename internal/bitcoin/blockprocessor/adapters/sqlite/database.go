package sqlite

import (
	"os"
	"path/filepath"

	"github.com/iangregsondev/deblockprocessor/internal/bitcoin/blockprocessor/entities"
	oswrapper "github.com/iangregsondev/deblockprocessor/internal/wrappers/os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseProvider is a struct that holds the database connection.
type DatabaseProvider struct {
	osWrapper    oswrapper.OS
	DB           *gorm.DB
	databaseFile string
}

// NewDatabaseProvider returns an instance of DatabaseProvider implementing the Database interface.
func NewDatabaseProvider(osWrapper oswrapper.OS, databaseFile string) Database {
	return &DatabaseProvider{
		osWrapper:    osWrapper,
		databaseFile: databaseFile,
	}
}

// Connect initializes the SQLite database connection.
func (s *DatabaseProvider) Connect() error {
	var err error

	// Ensure the directory for the database file exists
	dir := filepath.Dir(s.databaseFile)
	if err := s.osWrapper.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	s.DB, err = gorm.Open(
		sqlite.Open(s.databaseFile), &gorm.Config{
			// TODO: Get this from the config
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the database connection.
func (s *DatabaseProvider) GetDB() *gorm.DB {
	return s.DB
}

// Migrate performs the database schema migration for the Block model.
func (s *DatabaseProvider) Migrate() error {
	return s.DB.AutoMigrate(&entities.BlockHeight{})
}
