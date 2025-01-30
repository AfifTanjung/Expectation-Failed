package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RawRepository interface {
	CreateByAlatId(ctx context.Context, raw model.Raw) error
	GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error)
}

type rawRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRawRepository(db *gorm.DB, logger *log.Logger) RawRepository {
	return &rawRepository{db: db, logger: logger}
}

func (r *rawRepository) CreateByAlatId(ctx context.Context, raw model.Raw) error {
	r.logger.Printf("[INFO] Creating Raw Data for AlatID: %s", raw.AlatID)
	var alat model.Alat
	if err := r.db.WithContext(ctx).First(&alat, "id = ?", raw.AlatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("alat with ID %s not found", raw.AlatID)
		}
		return fmt.Errorf("failed to validate AlatID: %w", err)
	}
	if err := r.db.WithContext(ctx).Create(&raw).Error; err != nil {
		r.logger.Printf("[ERROR] Error Creating Raw Data for AlatID %s: %v", raw.AlatID, err)
		return fmt.Errorf("failed to create raw data: %w", err)
	}
	r.logger.Printf("[INFO] Successfully Created Raw Data for AlatID: %s", raw.AlatID)
	return nil
}

func (r *rawRepository) GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error) {
	r.logger.Printf("[INFO] Fetching Alat and Raw Data for AlatID: %s", alatID)

	if alatID == uuid.Nil {
		return nil, nil, errors.New("alatID cannot be empty")
	}
	var alat model.Alat
	if err := r.db.WithContext(ctx).Preload("Raws").First(&alat, "id = ?", alatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("alat with ID %s not found", alatID)
		}
		return nil, nil, fmt.Errorf("failed to fetch alat: %w", err)
	}

	return &alat, alat.Raws, nil
}
