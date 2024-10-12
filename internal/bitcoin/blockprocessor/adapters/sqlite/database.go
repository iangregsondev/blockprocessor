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

type DBConnection struct {
	osWrapper    oswrapper.OS
	DB           *gorm.DB
	databaseFile string
}

// NewSqliteDatabase returns an instance of DBConnection implementing the Database interface.
func NewSqliteDatabase(osWrapper oswrapper.OS, databaseFile string) Database {
	return &DBConnection{
		osWrapper:    osWrapper,
		databaseFile: databaseFile,
	}
}

// Connect initializes the SQLite database connection.
func (s *DBConnection) Connect() error {
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
func (s *DBConnection) GetDB() *gorm.DB {
	return s.DB
}

// Migrate performs the database schema migration for the Block model.
func (s *DBConnection) Migrate() error {
	return s.DB.AutoMigrate(&entities.BlockHeight{})
}
