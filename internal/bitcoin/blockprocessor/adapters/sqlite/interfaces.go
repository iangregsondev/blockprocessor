package sqlite

import "gorm.io/gorm"

type Database interface {
	Connect() error
	Migrate() error
	GetDB() *gorm.DB
}
