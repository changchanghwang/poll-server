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

type KakaoOauthClient struct {
	httpClient *http.Client
}

func newKakaoClient() *KakaoOauthClient {
	return &KakaoOauthClient{
		httpClient: http.DefaultClient,
	}
}

func (kakao *KakaoOauthClient) GetToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", config.Oauth.Kakao.ClientId)
	data.Set("client_secret", config.Oauth.Kakao.ClientSecret)
	data.Set("redirect_uri", config.Oauth.Kakao.RedirectUri)
	data.Set("code", code)

	res, err := kakao.httpClient.Post("https://kauth.kakao.com/oauth/token", "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
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

func (kakao *KakaoOauthClient) GetUserInfo(accessToken string) (*OauthUserInfo, error) {
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := kakao.httpClient.Do(req)
	if err != nil {
		return nil, httpError.New(httpCode.Unauthorized, err.Error(), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, httpError.New(httpCode.Unauthorized, string(body), "인증에 실패했습니다. 잠시 후 다시 시도해주세요.")
	}

	userInfo, err := kakao.parseUserInfo(body)
	if err != nil {
		return nil, httpError.Wrap(err)
	}

	return userInfo, nil
}

func (kakao *KakaoOauthClient) parseUserInfo(body []byte) (*OauthUserInfo, error) {
	var result struct {
		KakaoAccount struct {
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
		Email:           result.KakaoAccount.Email,
		Nickname:        result.KakaoAccount.Profile.Nickname,
		ProfileImageUrl: result.KakaoAccount.Profile.ProfileImageUrl,
	}

	return &userInfo, nil
}
