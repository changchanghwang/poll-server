package ddd

import (
	"time"

	"gorm.io/gorm"
)

type Aggregate struct {
	CreatedAt time.Time `json:"-" gorm:"column:created_at;autoCreateTime:nano"`
	UpdatedAt time.Time `json:"-" gorm:"column:updated_at;autoUpdateTime:nano"`
}

type SoftDeletableAggregate struct {
	Aggregate
	DeletedAt *gorm.DeletedAt `json:"-" gorm:"column:deleted_at"`
}
