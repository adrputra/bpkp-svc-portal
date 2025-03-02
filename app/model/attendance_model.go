package model

import "time"

type RequestUserAttendances struct {
	Username      string `json:"username" validate:"required"`
	InstitutionID string `json:"institution_id"`
	RoleID        string `json:"role_id"`
	RoleLevel     int    `json:"role_level"`
	Filter        Filter `json:"filter"`
}

type Attendance struct {
	ID        string    `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username" validate:"required"`
	CheckIn   time.Time `json:"check_in" gorm:"column:check_in"`
	CheckOut  time.Time `json:"check_out" gorm:"column:check_out"`
	StatusIn  string    `json:"status_in" gorm:"column:status_in"`
	StatusOut string    `json:"status_out" gorm:"column:status_out"`
	RemarkIn  string    `json:"remark_in" gorm:"column:remark_in"`
	RemarkOut string    `json:"remark_out" gorm:"column:remark_out"`
	SourceIn  string    `json:"source_in" gorm:"column:source_in"`
	SourceOut string    `json:"source_out" gorm:"column:source_out"`
}

type UserAttendance struct {
	ID          string    `json:"id" gorm:"column:id"`
	Username    string    `json:"username" gorm:"column:username" validate:"required"`
	CheckIn     time.Time `json:"check_in" gorm:"column:check_in"`
	CheckOut    time.Time `json:"check_out" gorm:"column:check_out"`
	StatusIn    string    `json:"status_in" gorm:"column:status_in"`
	StatusOut   string    `json:"status_out" gorm:"column:status_out"`
	RemarkIn    string    `json:"remark_in" gorm:"column:remark_in"`
	RemarkOut   string    `json:"remark_out" gorm:"column:remark_out"`
	SourceIn    string    `json:"source_in" gorm:"column:source_in"`
	SourceOut   string    `json:"source_out" gorm:"column:source_out"`
	Fullname    string    `json:"fullname" gorm:"column:fullname"`
	Shortname   string    `json:"shortname" gorm:"column:shortname"`
	Email       string    `json:"email" gorm:"column:email"`
	Gender      string    `json:"gender" gorm:"column:gender"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
}
