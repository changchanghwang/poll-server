package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	FindByUserId(db *gorm.DB, userId uuid.UUID) (*RefreshToken, bool, error)
	FindOneOrFail(db *gorm.DB, id uuid.UUID) (*RefreshToken, error)
	Save(db *gorm.DB, refreshToken *RefreshToken) error
}
