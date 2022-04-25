package models

import (
	_ "github.com/jinzhu/gorm"
)

type Homework struct {
	Id       int     `json:"id" gorm:"primary_key"`
	Grade    string  `json:"grade"`
	Class    string  `json:"class"`
	Course   string  `json:"course"`
	Task     string  `json:"task"`
	Reported *string `json:"reportedBy" default:"null"`
	Author   string  `json:"author"`
	Deadline int64   `json:"deadline"` // to have a unixTimestamp. Gorm and go have bad timestamp support.
}
