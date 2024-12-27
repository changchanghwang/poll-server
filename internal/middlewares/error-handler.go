package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
	httpError "poll.ant/internal/libs/http/http-error"
	httpResponse "poll.ant/internal/libs/http/http-response"
)

func ErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	if err != nil {
		e := httpError.UnWrap(err)

		//TODO: log error with something (e.g. Sentry, ELK, File, etc.)
		fmt.Println(e.Stack)

		ctx.JSON(e.Code, httpResponse.Response{Data: e.ClientMessage})
		return
	}
}
