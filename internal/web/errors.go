package web

import (
	"fmt"
	"net/http"

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

func FancyHandleNotFound(c echo.Context) error {
	// 🐱
	c.Response().WriteHeader(http.StatusNotFound)
	_, err := c.Response().Write([]byte(
		`<html><body><center><img src="https://http.cat/404" alt="Not Found" /></center></body></html>`,
	))
	return err
}
