package database

import (
	"fmt"

	"github.com/iangregsondev/deblockprocessor/internal/ethereum/blockprocessor/adapters/sqlite"
)

type Service struct {
	dbAdapter sqlite.Database
}

// NewService initializes the service with the given database implementation.
func NewService(dbAdapter sqlite.Database) *Service {
	return &Service{dbAdapter: dbAdapter}
}

// Setup initializes the database connection and runs migrations.
func (s *Service) Setup() error {
	err := s.dbAdapter.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = s.dbAdapter.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate the database: %w", err)
	}

	return nil
}
