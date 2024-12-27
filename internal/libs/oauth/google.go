package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"poll.ant/internal/config"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

type GoogleOauthClient struct {
	httpClient *http.Client
}

func newGoogleClient() *GoogleOauthClient {
	return &GoogleOauthClient{
		httpClient: http.DefaultClient,
	}
}

func (google *GoogleOauthClient) GetToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", config.Oauth.Google.ClientId)
	data.Set("client_secret", config.Oauth.Google.ClientSecret)
	data.Set("redirect_uri", config.Oauth.Google.RedirectUri)
	data.Set("code", code)

	res, err := google.httpClient.Post("https://nid.google.com/oauth2.0/token", "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
		}

		return "", httpError.New(httpCode.Unauthorized, string(body), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	return result.AccessToken, nil
}

func (google *GoogleOauthClient) GetUserInfo(accessToken string) (*OauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://openapi.google.com/v1/nid/me", nil)
	if err != nil {
		return nil, httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := google.httpClient.Do(req)
	if err != nil {
		return nil, httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, httpError.New(httpCode.Unauthorized, string(body), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	userInfo, err := google.parseUserInfo(body)
	if err != nil {
		return nil, httpError.Wrap(err)
	}

	return userInfo, nil
}

func (google *GoogleOauthClient) parseUserInfo(body []byte) (*OauthUserInfo, error) {
	var result struct {
		GoogleAccount struct {
			Profile struct {
				Nickname        string `json:"nickname"`
				ProfileImageUrl string `json:"profile_image_url"`
			} `json:"profile"`
			Email string `json:"email"`
		} `json:"kakao_account"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	userInfo := OauthUserInfo{
		Email:           result.GoogleAccount.Email,
		Nickname:        result.GoogleAccount.Profile.Nickname,
		ProfileImageUrl: result.GoogleAccount.Profile.ProfileImageUrl,
	}

	return &userInfo, nil
}
