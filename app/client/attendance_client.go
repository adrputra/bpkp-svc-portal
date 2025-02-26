package client

import (
	"context"
	"face-recognition-svc/app/model"
	"face-recognition-svc/app/utils"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type InterfaceAttendanceClient interface {
	GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.Attendance, error)
	GetTodayAttendances(ctx context.Context, username string) (*model.UserAttendance, error)
	CheckIn(ctx context.Context, request *model.Attendance) error
}

type AttendanceClient struct {
	db *gorm.DB
}

func NewAttendanceClient(db *gorm.DB) *AttendanceClient {
	return &AttendanceClient{db: db}
}

func (c *AttendanceClient) GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.Attendance, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetUserAttendances")
	defer span.Finish()

	var response []*model.Attendance

	sb := strings.Builder{}
	if request.Filter.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", request.Filter.Limit))
	}

	if request.Filter.SortType != "" {
		sb.WriteString(fmt.Sprintf(" %s", request.Filter.SortType))
	}

	query := "SELECT * FROM attendance WHERE username = ?"

	utils.LogEvent(span, "Query", query+sb.String())

	err := c.db.Debug().Raw(query, request.Username, sb.String()).Scan(&response).Error

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return response, nil
}

func (c *AttendanceClient) GetTodayAttendances(ctx context.Context, username string) (*model.UserAttendance, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetTodayAttendances")
	defer span.Finish()

	utils.LogEvent(span, "Request", username)

	var response *model.UserAttendance

	query := "SELECT a.*, u.fullname, u.shortname, u.email, u.gender  FROM attendance AS a INNER JOIN users AS u ON a.username = u.username WHERE a.username = ? AND DATE(a.check_in) = CURDATE()"
	utils.LogEvent(span, "Query", query)

	err := c.db.Debug().Raw(query, username).Scan(&response).Error

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	return response, nil
}

func (c *AttendanceClient) CheckIn(ctx context.Context, request *model.Attendance) error {
	span, _ := utils.SpanFromContext(ctx, "Client: CheckIn")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	var args []interface{}

	args = append(args, request.Username, request.CheckIn, nil, request.StatusIn, request.StatusOut, request.RemarkIn, request.RemarkOut)
	query := "INSERT INTO attendance (username, check_in, check_out, status_in, status_out, remark_in, remark_out) VALUES (?, ?, ?, ?, ?, ?, ?)"

	err := c.db.Debug().Exec(query, args...).Error

	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}
