package web

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// WrapHTTPError wraps error with message into echo.HTTPError
func WrapHTTPError(code int, err error, msg string) *echo.HTTPError {
	return echo.NewHTTPError(code, fmt.Sprintf("%s: %s", msg, err)).SetInternal(err)
}

// ToHTTPError creates echo.HTTPError from error
func ToHTTPError(code int, err error) *echo.HTTPError {
	return echo.NewHTTPError(code, err.Error())
}

// FormatHTTPError returns echo.HTTPError with formatted message
func FormatHTTPError(code int, msg string, args ...interface{}) *echo.HTTPError {
	return echo.NewHTTPError(code, fmt.Sprintf(msg, args...))
}
