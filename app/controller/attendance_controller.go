package controller

import (
	"context"
	"face-recognition-svc/app/client"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"time"
)

type InterfaceAttendanceController interface {
	GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.Attendance, error)
	GetTodayAttendances(ctx context.Context) (*model.UserAttendance, error)
	CheckIn(ctx context.Context, request *model.Attendance) error
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

func (uc *AttendanceController) GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.Attendance, error) {
	span, ctx := utils.SpanFromContext(ctx, "Controller: GetUserAttendances")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

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

	session, err := utils.GetMetadata(ctx)
	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	request.Username = session.Username
	request.CheckIn = time.Now()

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

	targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), parsedTime.Hour(), parsedTime.Minute(), 0, 0, time.Local)
	if time.Now().Compare(targetTime) == -1 {
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
