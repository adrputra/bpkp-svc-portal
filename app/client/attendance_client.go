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
	GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.UserAttendance, error)
	GetTodayAttendances(ctx context.Context, username string) (*model.UserAttendance, error)
	CheckIn(ctx context.Context, request *model.Attendance) error
	CheckOut(ctx context.Context, request *model.Attendance) error
}

type AttendanceClient struct {
	db *gorm.DB
}

func NewAttendanceClient(db *gorm.DB) *AttendanceClient {
	return &AttendanceClient{db: db}
}

func (c *AttendanceClient) GetUserAttendances(ctx context.Context, request *model.RequestUserAttendances) ([]*model.UserAttendance, error) {
	span, _ := utils.SpanFromContext(ctx, "Client: GetUserAttendances")
	defer span.Finish()

	var response []*model.UserAttendance

	sb := strings.Builder{}

	if request.RoleLevel == 3 {
		sb.WriteString(fmt.Sprintf(" AND u.username = %s", request.Username))
	}

	if request.RoleLevel == 2 {
		sb.WriteString(fmt.Sprintf(" AND u.institution_id = %s", request.InstitutionID))
	}

	if request.Filter.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", request.Filter.Limit))
	}

	if request.Filter.SortType != "" {
		sb.WriteString(fmt.Sprintf(" %s", request.Filter.SortType))
	}

	query := "SELECT a.*, u.fullname, u.shortname, u.email, u.gender  FROM attendance AS a INNER JOIN users AS u ON a.username = u.username"

	utils.LogEvent(span, "Query", query+sb.String())

	err := c.db.Debug().Raw(query + sb.String()).Scan(&response).Error

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

	args = append(args, request.Username, request.CheckIn, request.StatusIn, request.RemarkIn, request.SourceIn, request.Username)
	query := "INSERT INTO attendance (username, check_in, status_in, remark_in, source_in) SELECT ?, ?, ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM attendance WHERE username = ? AND DATE(check_in) = CURDATE())"

	err := c.db.Debug().Exec(query, args...).Error

	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}

func (c *AttendanceClient) CheckOut(ctx context.Context, request *model.Attendance) error {
	span, _ := utils.SpanFromContext(ctx, "Client: CheckOut")
	defer span.Finish()

	utils.LogEvent(span, "Request", request)

	var args []interface{}

	args = append(args, request.CheckOut, request.StatusOut, request.RemarkOut, request.SourceOut, request.Username)
	query := "UPDATE attendance SET check_out = ?, status_out = ?, remark_out = ?, source_out = ? WHERE username = ? AND DATE(check_in) = CURDATE()"

	err := c.db.Debug().Exec(query, args...).Error

	if err != nil {
		utils.LogEventError(span, err)
		return err
	}

	return nil
}
