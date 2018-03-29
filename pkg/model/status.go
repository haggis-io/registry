package model

import (
	"database/sql/driver"
)

const (
	StatusPENDING  Status = "PENDING"
	StatusAPPROVED Status = "APPROVED"
	StatusDECLINED Status = "DECLINED"
)

type Status string

func (s *Status) Scan(value interface{}) error { *s = Status(value.([]byte)); return nil }
func (s Status) Value() (driver.Value, error)  { return string(s), nil }

func (s *Status) Pending() {
	*s = StatusPENDING
}

func (s *Status) Approve() {
	*s = StatusAPPROVED
}

func (s *Status) Decline() {
	*s = StatusDECLINED
}
