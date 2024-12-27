package domain

import (
	"github.com/google/uuid"
	"poll.ant/internal/libs/ddd"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

type User struct {
	ddd.SoftDeletableAggregate
	Id         uuid.UUID `json:"id" gorm:"primaryKey; type:uuid"`
	Email      string    `json:"email" gorm:"unique;type:varchar(50); not null;"`
	Name       string    `json:"name" gorm:"type:varchar(50); not null;"`
	Provider   string    `json:"provider" gorm:"type:varchar(50); not null;"`
	ProviderId string    `json:"provider_id" gorm:"type:varchar(50); column:provider_id; not null;"`
	Role       string    `json:"role" gorm:"type:varchar(50); not null;"`
}

func (u *User) TableName() string {
	return "user"
}

func New(email, name, provider string) (*User, error) {
	uuId, err := uuid.NewV7()
	if err != nil {
		return nil, httpError.New(httpCode.InternalServerError, "Failed to create new user. Can not generate uuid.", "")
	}

	user := &User{
		Id:       uuId,
		Email:    email,
		Name:     name,
		Provider: provider,
	}

	return user, nil
}

type UpdateType struct {
	Name *string `json:"name,omitempty"`
}

func (u *User) Update(update UpdateType) {
	if update.Name != nil {
		u.Name = *update.Name
	}
}
