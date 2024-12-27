package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"poll.ant/internal/config"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

func Sign(userId any, expiredAfter time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", httpError.New(httpCode.InternalServerError, "Failed to encode access token. Can not convert claims to jwt.MapClaims.", "")
	}

	if userId != nil {
		claims["userId"] = userId
	}
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(expiredAfter).Unix()

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to encode access token. %v", err), "")
	}

	return tokenString, nil
}
