package model

import (
	"time"

	"github.com/google/uuid"
)

type Maintenance struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AlatID    uuid.UUID `gorm:"type:uuid;not null" json:"alat_id"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Nama      string    `json:"nama"`
	Laporan   string    `json:"laporan"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}
