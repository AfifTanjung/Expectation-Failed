package model

import (
	"time"
)

type Valid struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	RawID      uint      `gorm:"not null;unique"`
	Co2        string    `gorm:"type:varchar(50)" json:"co2"`
	Pm2_5      string    `gorm:"type:varchar(50)" json:"pm2_5"`
	Temp       string    `gorm:"type:varchar(50)" json:"temp"`
	Hum        string    `gorm:"type:varchar(50)" json:"hum"`
	Ws         string    `gorm:"type:varchar(50)" json:"ws"` // Wind Speed
	Wd         string    `gorm:"type:varchar(50)" json:"wd"` // Wind Direction
	WaktuMasuk time.Time `json:"waktuMasuk"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
