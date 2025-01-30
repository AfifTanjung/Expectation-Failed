package model

import (
	"time"

	"github.com/google/uuid"
)

type Raw struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	AlatID           uuid.UUID `gorm:"type:uuid;not null" json:"alatId"`
	Co2              string    `gorm:"type:varchar(50)" json:"co2"`
	Pm2_5            string    `gorm:"type:varchar(50)" json:"pm2_5"`
	Temp             string    `gorm:"type:varchar(50)" json:"temp"`
	Hum              string    `gorm:"type:varchar(50)" json:"hum"`
	Ws               string    `gorm:"type:varchar(50)" json:"ws"` // Wind Speed
	Wd               string    `gorm:"type:varchar(50)" json:"wd"` // Wind Direction
	WaktuMasuk       time.Time `json:"waktuMasuk"`
	IsValidated      bool      `gorm:"default:false" json:"is_validated"`
	ValidationStatus string    `gorm:"type:varchar(50);default:'pending'" json:"validation_status"`
	Valid            Valid     `gorm:"foreignKey:RawID;constraint:OnDelete:CASCADE;"`
	Alat             Alat      `gorm:"foreignKey:AlatID;constraint:OnDelete:CASCADE;"`
}
