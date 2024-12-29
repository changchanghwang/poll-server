package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
)

type oauthClient interface {
	parseUserInfo([]byte) (*OauthInfo, error)
	getConfig() oauthConfig
}
type OAuthProvider struct {
	clients    map[string]oauthClient
	httpClient *http.Client
}

func NewOAuthProvider() *OAuthProvider {
	return &OAuthProvider{
		clients: map[string]oauthClient{
			"GOOGLE": newGoogleClient(),
		},
		httpClient: &http.Client{},
	}
}

func (o *OAuthProvider) client(provider string) (oauthClient, error) {
	client, ok := o.clients[provider]
	if !ok {
		return nil, httpError.New(httpCode.BadRequest, fmt.Sprintf("unknown provider: %s", provider), "")
	}
	return client, nil
}

func (o *OAuthProvider) OAuth(provider, code string) (*OauthInfo, error) {
	client, err := o.client(provider)
	if err != nil {
		return nil, err
	}
	// token 가져오기
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", client.getConfig().clientId)
	data.Set("client_secret", client.getConfig().clientSecret)
	data.Set("redirect_uri", client.getConfig().redirectUri)
	data.Set("code", code)

	res, err := o.httpClient.Post(client.getConfig().tokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to request kakao oauth token, %v", err), "")
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, httpError.New(httpCode.Status{Code: res.StatusCode}, fmt.Sprintf("unexpected status code: %d, body: %s", res.StatusCode, string(body)), "")
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(res.Body).Decode(&result)

	// user 정보 가져오기
	req, err := http.NewRequest("GET", client.getConfig().userInfoURL, nil)
	if err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to create request, %v", err), "")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", result.AccessToken))

	res, err = o.httpClient.Do(req)

	if err != nil {
		return nil, httpError.New(httpCode.InternalServerError, fmt.Sprintf("Failed to request user info, %v", err), "")
	}
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, httpError.New(httpCode.Status{Code: res.StatusCode}, fmt.Sprintf("unexpected status code: %d, body: %s", res.StatusCode, string(body)), "Internal Server Error")
	}

	return client.parseUserInfo(body)
}
