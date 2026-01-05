package response

import (
	"net/http"

	"github.com/kalom60/cashflow/internal/constant/errors"
	"github.com/labstack/echo/v4"

	"github.com/joomcode/errorx"
)

func SendSuccessResponse(c echo.Context, statusCode int, data any) error {
	response := SuccessResponse{
		Data: data,
	}

	return c.JSON(statusCode, response)
}

func SendErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorResponseFormat{
		Message: message,
	})
}

func SendErrorResponseFormated(c echo.Context, err error) error {
	statusCode := http.StatusInternalServerError
	message := "Internal Server Error"
	if err != nil {
		message = err.Error()
		for _, e := range errors.Error {
			if errorx.IsOfType(err, e.Type) {
				statusCode = e.StatusCode
				er := errorx.Cast(err)
				message = er.Message()
				break
			}
		}
	}

	response := ErrorResponseFormat{
		Message: message,
	}

	return c.JSON(statusCode, response)
}

func GetErrorFrom(err error) *ErrorResponse {
	for _, e := range errors.Error {
		if errorx.IsOfType(err, e.Type) {
			er := errorx.Cast(err)
			res := ErrorResponse{
				Message: er.Message(),
			}

			return &res
		}
	}

	return &ErrorResponse{
		Message: "Unknown server error",
	}
}
