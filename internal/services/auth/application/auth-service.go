package application

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"poll.ant/internal/config"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
	"poll.ant/internal/libs/jwt"
	"poll.ant/internal/libs/oauth"
	"poll.ant/internal/services/auth/domain"
	"poll.ant/internal/services/auth/dto"
	users "poll.ant/internal/services/users/domain"
)

type AuthService struct {
	refreshTokenRepository domain.RefreshTokenRepository
	userRepository         users.UserRepository
	oauthProvider          *oauth.OAuthProvider
	db                     *gorm.DB
}

func NewAuthService(
	refreshTokenRepository domain.RefreshTokenRepository,
	userRepository users.UserRepository,
	oauthProvider *oauth.OAuthProvider,
	db *gorm.DB,
) *AuthService {
	return &AuthService{
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
		oauthProvider:          oauthProvider,
		db:                     db,
	}
}

func (s *AuthService) OAuth(code, provider string) (*dto.OauthResponse, error) {
	var (
		responseDto *dto.OauthResponse
		err         error
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		userInfo, err := s.oauthProvider.OAuth(provider, code)
		if err != nil {
			return httpError.Wrap(err)
		}

		user, exist, err := s.userRepository.FindByEmail(tx, userInfo.Email)
		if err != nil {
			return httpError.Wrap(err)
		}

		if !exist {
			newUser, err := users.New(userInfo.Email, userInfo.Name, provider)
			if err != nil {
				return httpError.Wrap(err)
			}

			err = s.userRepository.Save(tx, newUser)
			if err != nil {
				return httpError.Wrap(err)
			}

			user = newUser
		}

		accessToken, refreshToken, err := s.createToken(tx, user.Id, provider)
		if err != nil {
			return httpError.Wrap(err)
		}

		responseDto = &dto.OauthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return responseDto, nil
}

func (service *AuthService) createToken(tx *gorm.DB, userId uuid.UUID, clientInfo string) (string, string, error) {
	accessTokenExpiredAfterHour, err := strconv.Atoi(config.AccessTokenExpiredAfterHour)
	if err != nil {
		return "", "", httpError.New(httpCode.InternalServerError, "Failed to convert access token expired after hour to integer.", "")
	}

	refreshTokenExpiredAfterHour, err := strconv.Atoi(config.RefreshTokenExpiredAfterHour)
	if err != nil {
		return "", "", httpError.New(httpCode.InternalServerError, "Failed to convert refresh token expired after hour to integer.", "")
	}

	accessToken, err := jwt.Sign(userId, time.Hour*time.Duration(accessTokenExpiredAfterHour))
	if err != nil {
		return "", "", httpError.Wrap(err)
	}

	refreshToken, err := jwt.Sign(nil, time.Hour*time.Duration(refreshTokenExpiredAfterHour))
	if err != nil {
		return "", "", httpError.Wrap(err)
	}

	fmt.Println("###", userId)
	existRefreshToken, exist, err := service.refreshTokenRepository.FindByUserId(tx, userId)
	fmt.Println("$$$", existRefreshToken, exist)
	if err != nil {
		return "", "", httpError.Wrap(err)
	}

	if !exist {
		refreshTokenModel, err := domain.New(refreshToken, userId)
		fmt.Println("!!!", refreshTokenModel)
		if err != nil {
			return "", "", httpError.Wrap(err)
		}

		err = service.refreshTokenRepository.Save(tx, refreshTokenModel)

		if err != nil {
			return "", "", httpError.Wrap(err)
		}
	} else {
		existRefreshToken.Update(refreshToken)
		fmt.Println("@@@", existRefreshToken)
		err = service.refreshTokenRepository.Save(tx, existRefreshToken)
		if err != nil {
			return "", "", httpError.Wrap(err)
		}
	}

	return accessToken, refreshToken, nil
}
