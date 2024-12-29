package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(db *gorm.DB, email string) (*User, bool, error)
	FindOneOrFail(db *gorm.DB, id uuid.UUID) (*User, error)
	Save(db *gorm.DB, user *User) error
}
