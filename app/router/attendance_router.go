package router

import "github.com/labstack/echo/v4"

func InitAttendanceRoute(prefix string, e *echo.Group) {
	route := e.Group(prefix)
	service := factory.Service.attendance

	route.GET("", service.GetTodayAttendances)
	route.POST("", service.GetUserAttendances)
	route.POST("/checkin", service.CheckIn)
	route.POST("/checkout", service.CheckOut)
}
