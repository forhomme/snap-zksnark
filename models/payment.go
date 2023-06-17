package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	Id             string         `json:"id" gorm:"primary_key"`
	PartnerId      string         `json:"partnerId" gorm:"column:partner_id"`
	ConsumerId     string         `json:"consumerId" gorm:"column:consumer_id"`
	Amount         string         `json:"amount"`
	Currency       string         `json:"currency"`
	AdditionalInfo JSONB          `json:"additionalInfo" gorm:"column:additional_info"`
	CreatedAt      *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time     `json:"updatedAt,omitempty"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt,omitempty" sql:"index"`
}

// JSONB Interface for JSONB Field of yourTableName Table
type JSONB map[string]interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
