package repository

import (
	"context"
	"errors"
	"fmt"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RawRepository interface {
	CreateByAlatId(ctx context.Context, raw model.Raw) error
	GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error)
}

type rawRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewRawRepository(db *gorm.DB, logger *logrus.Logger) RawRepository {
	return &rawRepository{db: db, logger: logger}
}

func (r *rawRepository) CreateByAlatId(ctx context.Context, raw model.Raw) error {
	r.logger.WithField("AlatID", raw.AlatID).Info("[CreateByAlatId] Creating Raw Data")

	var alat model.Alat
	if err := r.db.WithContext(ctx).First(&alat, "id = ?", raw.AlatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithField("AlatID", raw.AlatID).Error("[CreateByAlatId] Alat not found")
			return fmt.Errorf("alat with ID %s not found", raw.AlatID)
		}
		r.logger.WithField("AlatID", raw.AlatID).WithError(err).Error("[CreateByAlatId] Error validating AlatID")
		return fmt.Errorf("failed to validate AlatID: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&raw).Error; err != nil {
		r.logger.WithField("AlatID", raw.AlatID).WithError(err).Error("[CreateByAlatId] Error creating raw data")
		return fmt.Errorf("failed to create raw data: %w", err)
	}

	r.logger.WithField("AlatID", raw.AlatID).Info("[CreateByAlatId] Successfully created raw data")
	return nil
}

func (r *rawRepository) GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error) {
	r.logger.WithField("AlatID", alatID).Info("[GetById] Fetching Alat and Raw Data")

	if alatID == uuid.Nil {
		r.logger.WithField("AlatID", alatID).Error("[GetById] AlatID cannot be empty")
		return nil, nil, errors.New("alatID cannot be empty")
	}

	var alat model.Alat
	if err := r.db.WithContext(ctx).Preload("Raws").First(&alat, "id = ?", alatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.WithField("AlatID", alatID).Error("[GetById] Alat not found")
			return nil, nil, fmt.Errorf("alat with ID %s not found", alatID)
		}
		r.logger.WithField("AlatID", alatID).WithError(err).Error("[GetById] Error fetching alat")
		return nil, nil, fmt.Errorf("failed to fetch alat: %w", err)
	}

	r.logger.WithField("AlatID", alatID).Info("[GetById] Successfully fetched Alat and Raw Data")
	return &alat, alat.Raws, nil
}
