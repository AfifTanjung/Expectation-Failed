package repository

import (
	"context"
	"log"
	"telkom-property-service/src/model"

	"gorm.io/gorm"
)

type ValidationRepository interface {
	GetPendingRawData(ctx context.Context) ([]model.Raw, error)
	UpdateRawStatus(ctx context.Context, raw model.Raw) error
	InsertValidData(ctx context.Context, valid model.Valid) error
}

type validationRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewValidationRepository(db *gorm.DB, logger *log.Logger) ValidationRepository {
	return &validationRepository{db: db, logger: logger}
}

func (r *validationRepository) GetPendingRawData(ctx context.Context) ([]model.Raw, error) {
	var raws []model.Raw
	r.logger.Println("[PRELOADING] data for validation")
	if err := r.db.WithContext(ctx).Preload("Alat").Where("is_validated = ?", false).Find(&raws).Error; err != nil {
		r.logger.Printf("[ERROR] fetching pending raw data: %v", err)
		return nil, err
	}
	return raws, nil
}

func (r *validationRepository) UpdateRawStatus(ctx context.Context, raw model.Raw) error {
	result := r.db.WithContext(ctx).Save(&raw)
	r.logger.Printf("[UPDATE] Status Data from Raw to Validation for ID: %d", raw.ID)
	if result.Error != nil {
		r.logger.Printf("[ERROR] updating raw data status: %v", result.Error)
		return result.Error
	}
	return nil
}

func (r *validationRepository) InsertValidData(ctx context.Context, valid model.Valid) error {
	result := r.db.WithContext(ctx).Create(&valid)
	r.logger.Printf("[VALIDATION] Inserting Validation Data for RawID: %d", valid.RawID)
	if result.Error != nil {
		r.logger.Printf("[ERROR] inserting validation data: %v", result.Error)
		return result.Error
	}
	r.logger.Printf("[SUCCESS] Validation data inserted for RawID: %d", valid.RawID)
	return nil
}
