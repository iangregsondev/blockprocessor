package entities

type BlockHeight struct {
	ID          uint  `gorm:"primaryKey"`
	BlockHeight int64 `gorm:"unique;not null"`
}
