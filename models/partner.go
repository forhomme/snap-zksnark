package models

import (
	"gorm.io/gorm"
	"time"
)

type Partners struct {
	Id          string         `json:"id" gorm:"primary_key"`
	ReferenceNo string         `json:"referenceNo" gorm:"column:reference_no;uniqueIndex:partner_reference_no_uindex"`
	Username    string         `json:"username" gorm:"column:username;uniqueIndex:partner_username_uindex"`
	Password    string         `json:"password"`
	CreatedAt   *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" sql:"index"`
}

func (Partners) TableName() string {
	return "partners"
}
