package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	IAQMS = "IAQMS"
	OAQMS = "OAQMS"
)

type Alat struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Nama         string        `gorm:"type:varchar(100);not null" json:"nama"`
	Type         string        `gorm:"type:varchar(10);default:'IAQMS'" json:"type"`
	Lantai       int           `gorm:"not null" json:"lantai"`
	Stasiun      string        `gorm:"type:varchar(50);not null" json:"stasiun"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	Raws         []Raw         `gorm:"foreignKey:AlatID;constraint:OnDelete:CASCADE;" json:"raws"`
	Maintenances []Maintenance `gorm:"foreignKey:AlatID;constraint:OnDelete:CASCADE;" json:"maintenances"`
}
