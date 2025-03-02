package controller

import (
	"bpkp-svc-portal/app/client"
	"bpkp-svc-portal/app/model"
	"bpkp-svc-portal/app/utils"
	"context"
	"strings"
	"time"
)

type InterfaceAttendanceController interface {
	GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.UserAttendance, error)
	GetTodayAttendances(ctx context.Context) (*model.UserAttendance, error)
	CheckIn(ctx context.Context, request *model.Attendance) error
	CheckOut(ctx context.Context, request *model.Attendance) error
}

type AttendanceController struct {
	attendanceClient client.InterfaceAttendanceClient
	paramClient      client.InterfaceParamClient
}

func NewAttendanceController(attendanceClient client.InterfaceAttendanceClient, paramClient client.InterfaceParamClient) *AttendanceController {
	return &AttendanceController{
		attendanceClient: attendanceClient,
		paramClient:      paramClient,
	}
}

func (uc *AttendanceController) GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.UserAttendance, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetUserAttendances")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	roleLevel, err := uc.paramClient.GetParameterByKey(ctx, "role-level-1")
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}
	if utils.Contains(strings.Split(roleLevel.Value, ";"), request.RoleID) {
		request.RoleLevel = 1
	} else {
		roleLevel, err = uc.paramClient.GetParameterByKey(ctx, "role-level-2")
		if err != nil {
			utils.LogEventError(span, err)
			return nil, err
		}
		if utils.Contains(strings.Split(roleLevel.Value, ";"), request.RoleID) {
			request.RoleLevel = 2
		} else {
			request.RoleLevel = 3
		}
	}

	res, err := uc.attendanceClient.GetUserAttendances(ctx, request)

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return res, nil
}

func (uc *AttendanceController) GetTodayAttendances(ctx context.Context) (*model.UserAttendance, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetTodayAttendances")
	defer span.Finish()

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	username := session.Username
	utils.LogEvent(span, "Request", session.Username)

	res, err := uc.attendanceClient.GetTodayAttendances(ctx, username)
	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return res, nil
}

func (uc *AttendanceController) CheckIn(ctx context.Context, request *model.Attendance) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CheckIn")
	defer span.Finish()

	request.CheckIn = utils.LocalTime()

	checkInThreshold, err := uc.paramClient.GetParameterByKey(ctx, "checkin-time")
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	parsedTime, err := time.Parse("15:04", checkInThreshold.Value)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	targetTime := time.Date(utils.LocalTime().Year(), utils.LocalTime().Month(), utils.LocalTime().Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	if utils.LocalTime().Compare(targetTime) == -1 {
		request.StatusIn = "On Time"
	} else {
		request.StatusIn = "Late"
	}

	utils.LogEvent(span, "Request", request)

	err = uc.attendanceClient.CheckIn(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Check In")

	return nil
}

func (uc *AttendanceController) CheckOut(ctx context.Context, request *model.Attendance) error {
	span, ctx := utils.SpanFromContext(ctx, "Controller: CheckIn")
	defer span.Finish()

	request.CheckOut = utils.LocalTime()

	utils.LogEvent(span, "Request", request)
	checkInThreshold, err := uc.paramClient.GetParameterByKey(ctx, "checkout-time")
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	parsedTime, err := time.Parse("15:04", checkInThreshold.Value)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	targetTime := time.Date(utils.LocalTime().Year(), utils.LocalTime().Month(), utils.LocalTime().Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	if utils.LocalTime().Compare(targetTime) == -1 {
		request.StatusOut = "Early"
	} else {
		request.StatusOut = "Normal"
	}

	utils.LogEvent(span, "Request", request)

	err = uc.attendanceClient.CheckOut(ctx, request)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	utils.LogEvent(span, "Response", "Success Check Out")

	return nil
}
