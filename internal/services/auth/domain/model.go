package domain

import (
	"fmt"

	"github.com/google/uuid"
	"poll.ant/internal/libs/ddd"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

type RefreshToken struct {
	ddd.Aggregate
	Id     uuid.UUID `json:"id" gorm:"primaryKey;"`
	UserId uuid.UUID `json:"user_id" gorm:"type:varchar(36); not null; column:user_id;"`
	Token  string    `json:"token" gorm:"type:varchar(255); not null; column:token;"`
}

func (refreshToken *RefreshToken) TableName() string {
	return "refresh_token"
}

func (refreshToken *RefreshToken) Update(token string) {
	refreshToken.Token = token
}

func New(token string, userId uuid.UUID) (*RefreshToken, error) {
	fmt.Println("!@#$!@#$!@#$", userId)

	uuId, err := uuid.NewV7()
	if err != nil {
		return nil, httpError.New(httpCode.InternalServerError, "Failed to create new user. Can not generate uuid.", "")
	}

	return &RefreshToken{
		Id:     uuId,
		UserId: userId,
		Token:  token,
	}, nil
}
