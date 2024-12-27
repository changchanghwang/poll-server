package oauth

type OauthClient interface {
	GetToken(code string) (string, error)
	GetUserInfo(accessToken string) (*OauthUserInfo, error)
}

type OauthUserInfo struct {
	Email           string `json:"email"`
	Nickname        string `json:"nickname"`
	ProfileImageUrl string `json:"profileImageUrl"`
}
