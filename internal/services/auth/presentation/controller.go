package presentation

import (
	"github.com/gofiber/fiber/v2"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
	httpResponse "poll.ant/internal/libs/http/http-response"
	"poll.ant/internal/services/auth/application"
	"poll.ant/internal/services/auth/dto"
)

type AuthController struct {
	authService *application.AuthService
}

func NewAuthController(authService *application.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Route(r fiber.Router) {
	r.Post("/", c.oAuth)
}

// @Summary OAuth 인증
// @Description OAuth 인증을 수행한다.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.OauthRequestBody true "OAuth request body"
// @Success 200 {object} dto.OauthResponse "OAuth response"
// @Failure 400 {object} httpError.ErrorResponse "Bad Request"
// @Failure 401 {object} httpError.ErrorResponse "Unauthorized"
// @Failure 500 {object} httpError.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /users/auth [post]
func (c *AuthController) oAuth(ctx *fiber.Ctx) error {
	// 1. ctx destructuring
	// 2. dto validation
	var body dto.OauthRequestBody

	if err := ctx.BodyParser(&body); err != nil {
		return httpError.New(httpCode.BadRequest, "Invalid request body", "")
	}

	// if err := validate.ValidateDto(body); err != nil {
	// 	return httpError.Wrap(err)
	// }

	// 3. call application service method
	result, err := c.authService.OAuth(body.AuthorizationCode, body.Provider)
	if err != nil {
		return httpError.Wrap(err)
	}

	return ctx.Status(httpCode.Ok.Code).JSON(httpResponse.Response{Data: result})
}
