package routing

import (
	"github.com/kalom60/cashflow/platform/logger"
	"github.com/labstack/echo/v4"
)

type Route struct {
	Method     string
	Path       string
	Handler    echo.HandlerFunc
	Middleware []echo.MiddlewareFunc
}

func RegisterRoute(eg *echo.Group, routes []Route, log logger.Logger) {
	for _, route := range routes {
		eg.Add(route.Method, route.Path, route.Handler, route.Middleware...)
	}
}
