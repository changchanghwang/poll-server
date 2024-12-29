package oauth

type OauthInfo struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	ProfileImage string `json:"profileImage"`
	Provider     string `json:"provider"`
}

type oauthConfig struct {
	authURL      string
	tokenURL     string
	userInfoURL  string
	clientId     string
	clientSecret string
	redirectUri  string
	state        string
	scope        string
}
