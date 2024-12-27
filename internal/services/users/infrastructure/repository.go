package infrastructure

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
	"poll.ant/internal/services/users/domain"
)

type UserRepository interface {
	FindByEmail(db *gorm.DB, email string) (*domain.User, bool, error)
	FindOneOrFail(db *gorm.DB, id uuid.UUID) (*domain.User, error)
	Save(db *gorm.DB, user *domain.User) error
}

type UserRepositoryImpl struct {
	manager *gorm.DB
}

func NewUserRepository(manager *gorm.DB) UserRepository {
	return &UserRepositoryImpl{manager: manager}
}

func (r *UserRepositoryImpl) FindByEmail(db *gorm.DB, email string) (*domain.User, bool, error) {
	if db == nil {
		db = r.manager
	}

	users := []domain.User{}
	if err := db.Where("email = ?", email).Find(&users).Error; err != nil {
		return nil, false, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to findByEmail user. %v", err), "")
	}
	if len(users) == 0 {
		return nil, false, nil
	}
	return &users[0], true, nil
}

func (r *UserRepositoryImpl) FindOneOrFail(db *gorm.DB, id uuid.UUID) (*domain.User, error) {
	if db == nil {
		db = r.manager
	}

	var user *domain.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to findById user. %s", err.Error()), "")
	}
	if user == nil {
		return nil, httpError.New(httpCode.NotFound, fmt.Sprintf("User(%s) not found", id.String()), "")
	}

	return user, nil
}

func (r *UserRepositoryImpl) Save(db *gorm.DB, user *domain.User) error {
	if db == nil {
		db = r.manager
	}

	if err := db.Save(user).Error; err != nil {
		return httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to save user. %s", err.Error()), "")
	}
	return nil
}
