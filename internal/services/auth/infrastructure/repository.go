package infrastructure

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
	"poll.ant/internal/services/auth/domain"
)

type RefreshTokenRepositoryImpl struct {
	manager *gorm.DB
}

func New(manager *gorm.DB) domain.RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{manager: manager}
}

func (r *RefreshTokenRepositoryImpl) FindByUserId(db *gorm.DB, userId uuid.UUID) (*domain.RefreshToken, bool, error) {
	if db == nil {
		db = r.manager
	}

	refreshTokens := []domain.RefreshToken{}
	if err := db.Where("user_id = ?", userId).Find(&refreshTokens).Error; err != nil {
		return nil, false, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to findByUserId refreshToken. %v", err), "")
	}
	if len(refreshTokens) == 0 {
		return nil, false, nil
	}
	return &refreshTokens[0], true, nil
}

func (r *RefreshTokenRepositoryImpl) FindOneOrFail(db *gorm.DB, id uuid.UUID) (*domain.RefreshToken, error) {
	if db == nil {
		db = r.manager
	}

	var refreshToken *domain.RefreshToken
	if err := db.Where("id = ?", id).First(&refreshToken).Error; err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to findById refreshToken. %s", err.Error()), "")
	}
	if refreshToken == nil {
		return nil, httpError.New(httpCode.NotFound, fmt.Sprintf("RefreshToken(%s) not found", id.String()), "")
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) Save(db *gorm.DB, refreshToken *domain.RefreshToken) error {
	if db == nil {
		db = r.manager
	}

	if err := db.Save(refreshToken).Error; err != nil {
		return httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to save refreshToken. %s", err.Error()), "")
	}
	return nil
}
