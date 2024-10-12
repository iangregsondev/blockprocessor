package entities

type BlockNumber struct {
	ID          uint  `gorm:"primaryKey"`
	BlockNumber int64 `gorm:"unique;not null"`
}
