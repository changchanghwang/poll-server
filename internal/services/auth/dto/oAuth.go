package dto

type OauthResponse struct {
	AccessToken  string `json:"accessToken" validate:"required"`
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type OauthRequestBody struct {
	AuthorizationCode string `json:"authorizationCode" validate:"required"`
	Provider          string `json:"provider" validate:"required"`
}
