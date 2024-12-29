package oauth

import (
	"encoding/json"
	"fmt"

	"poll.ant/internal/config"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

type googleClient struct {
	config oauthConfig
}

func newGoogleClient() *googleClient {
	return &googleClient{config: oauthConfig{
		tokenURL:     "https://oauth2.googleapis.com/token",
		userInfoURL:  "https://oauth2.googleapis.com/tokeninfo",
		clientId:     config.Oauth.Google.ClientId,
		clientSecret: config.Oauth.Google.ClientSecret,
		redirectUri:  config.Oauth.Google.RedirectUri,
		state:        "ant",
		scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
	}}
}

func (c *googleClient) parseUserInfo(body []byte) (*OauthInfo, error) {
	var result struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("failed to decode user info response: %v", err), "")
	}

	userInfo := OauthInfo{
		Email:        result.Email,
		Name:         result.Name,
		ProfileImage: result.Picture,
	}

	return &userInfo, nil
}

func (c *googleClient) getConfig() oauthConfig {
	return c.config
}
