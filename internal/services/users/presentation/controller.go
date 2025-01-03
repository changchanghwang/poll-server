package presentation

import (
	"github.com/gofiber/fiber/v2"
	httpCode "poll.ant/internal/libs/http/http-code"
	httpError "poll.ant/internal/libs/http/http-error"
	httpResponse "poll.ant/internal/libs/http/http-response"
	"poll.ant/internal/libs/validate"
	"poll.ant/internal/services/users/application"
	"poll.ant/internal/services/users/domain"
	"poll.ant/internal/services/users/dto"
)

type UserController struct {
	userService *application.UserService
}

func NewUserController(userService *application.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) Route(r fiber.Router) {
	r.Patch("/", c.update)
}

// @Summary 사용자 정보 업데이트
// @Description 인증된 사용자의 정보를 업데이트한다.
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.UpdateUserRequestBody true "Updated user information"
// @Success 200 {object} dto.UpdateUserResponse "Updated user information"
// @Failure 400 {object} httpError.ErrorResponse "Bad Request"
// @Failure 401 {object} httpError.ErrorResponse "Unauthorized"
// @Failure 500 {object} httpError.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /users [patch]
func (c *UserController) update(ctx *fiber.Ctx) error {
	// 1. ctx destructuring
	user, ok := ctx.Locals("user").(*domain.User)
	if !ok {
		return httpError.New(httpCode.Unauthorized, "Unauthorized", "")
	}

	// 2. dto validation
	var body dto.UpdateUserRequestBody

	if err := ctx.BodyParser(&body); err != nil {
		return httpError.New(httpCode.BadRequest, "Invalid request body", "")
	}

	if err := validate.ValidateDto(body); err != nil {
		return httpError.Wrap(err)
	}

	// 3. call application service method
	result, err := c.userService.Update(user.Id, body)
	if err != nil {
		return httpError.Wrap(err)
	}

	return ctx.Status(httpCode.Ok.Code).JSON(httpResponse.Response{Data: result})
}
