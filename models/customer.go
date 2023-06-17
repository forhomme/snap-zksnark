package models

import (
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	Id         string         `json:"id" gorm:"primary_key"`
	KTP        string         `json:"ktp" gorm:"column:ktp;uniqueIndex:customers_ktp_uindex"`
	NoRek      string         `json:"account" gorm:"column:account;uniqueIndex:customers_account_uindex"`
	Name       string         `json:"name" gorm:"column:customer_name"`
	Branch     string         `json:"branch" gorm:"column:branch"`
	MotherName string         `json:"motherName" gorm:"column:mother_name"`
	CreatedAt  *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt,omitempty" sql:"index"`
}

func (Customer) TableName() string {
	return "customers"
}
