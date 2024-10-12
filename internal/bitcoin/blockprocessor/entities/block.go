package entities

type BlockHeight struct {
	ID     uint `gorm:"primaryKey"`
	Height int  `gorm:"unique;not null"`
}
