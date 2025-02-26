package service

import (
	"face-recognition-svc/app/controller"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type InterfaceAttendanceService interface {
	GetUserAttendances(e echo.Context) error
	GetTodayAttendances(e echo.Context) error
	CheckIn(e echo.Context) error
}

type AttendanceService struct {
	uc controller.InterfaceAttendanceController
}

func NewAttendanceService(uc controller.InterfaceAttendanceController) *AttendanceService {
	return &AttendanceService{uc: uc}
}

func (s *AttendanceService) GetUserAttendances(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetUserAttendances")
	defer span.Finish()

	var request *model.RequestUserAttendances

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	response, err := s.uc.GetUserAttendances(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get User Attendances",
		Data:    response,
	})
}

func (s *AttendanceService) CheckIn(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "CheckIn")
	defer span.Finish()

	var request *model.Attendance

	if err := e.Bind(&request); err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	utils.LogEvent(span, "Request", request)

	err := s.uc.CheckIn(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Check In",
		Data:    nil,
	})
}

func (s *AttendanceService) GetTodayAttendances(e echo.Context) error {
	ctx, span := utils.StartSpan(e, "GetTodayAttendances")
	defer span.Finish()

	res, err := s.uc.GetTodayAttendances(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return utils.LogError(e, err, nil)
	}

	return e.JSON(http.StatusOK, model.Response{
		Code:    200,
		Message: "Success Get Today Attendances",
		Data:    res,
	})
}
