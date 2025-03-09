package router

import (
	"bpkp-svc-portal/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InitPublicRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.user
	attendance := factory.Service.attendance

	route.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.Response{
			Code:    http.StatusOK,
			Message: "pong",
			Data:    nil,
		})
	})

	route.POST("/register", service.CreateNewUser)
	route.POST("/login", service.Login)
	route.GET("/metabase", service.EmbedMetabase)

	route.POST("/checkinout-rfid", attendance.CheckInOutRFID)
}
