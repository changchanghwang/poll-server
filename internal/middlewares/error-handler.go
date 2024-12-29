package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	httpError "poll.ant/internal/libs/http/http-error"
	httpResponse "poll.ant/internal/libs/http/http-response"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		return ctx.Status(e.Code).JSON(httpResponse.Response{Data: map[string]string{"errorMessage": e.Message}})
	}

	e := httpError.UnWrap(err)
	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	//TODO: log error with something (e.g. Sentry, ELK, File, etc.)
	fmt.Println(e.Stack)

	errResponse := httpResponse.Response{Data: map[string]string{"errorMessage": e.ClientMessage}}

	return ctx.Status(e.Code).JSON(errResponse)
}
