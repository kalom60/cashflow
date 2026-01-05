package handler

import "github.com/labstack/echo/v4"

type Payment interface {
	CreatePayment(c echo.Context) error
	GetPaymentDetails(c echo.Context) error
}
