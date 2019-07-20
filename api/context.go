package api

import (
	"github.com/labstack/echo"
	"strings"
)

type NovaApiContext struct {
	echo.Context
	// If address is not empty means this user is authenticated.
	Address string
}

func initNovaApiContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &NovaApiContext{c, ""}
		return next(cc)
	}
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*NovaApiContext)
		cc.Response().Header().Set(echo.HeaderServer, "Echo/3.0")

		novaAuthToken := cc.Request().Header.Get("Nova-Authentication")
		novaAuthTokens := strings.Split(novaAuthToken, "#")

		if len(novaAuthTokens) != 3 {
			return &ApiError{Code: -11, Desc: "Nova-Authentication should be like {address}#NOVA-AUTHENTICATION@{time}#{signature}"}
		}

		valid, err := nova.IsValidSignature(novaAuthTokens[0], novaAuthTokens[1], novaAuthTokens[2])
		if !valid || err != nil {
			return &ApiError{Code: -11, Desc: "Nova-Authentication valid failed, please check your authentication"}
		}
		cc.Address = strings.ToLower(novaAuthTokens[0])
		return next(cc)
	}
}
