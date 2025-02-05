package repository

import (
	"context"
	"telkom-property-service/src/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ValidationRepository interface {
	GetPendingRawData(ctx context.Context, limit, offset int) ([]model.Raw, error)
	UpdateRawStatus(ctx context.Context, rawID uint) error
	InsertValidData(ctx context.Context, valid model.Valid) error
	ValidateAndInsert(ctx context.Context, valid model.Valid) error
}

type validationRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

var _ ValidationRepository = (*validationRepository)(nil)

func NewValidationRepository(db *gorm.DB, logger *logrus.Logger) ValidationRepository {
	return &validationRepository{db: db, logger: logger}
}

func (r *validationRepository) GetPendingRawData(ctx context.Context, limit, offset int) ([]model.Raw, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var raws []model.Raw
	r.logger.WithFields(logrus.Fields{
		"action": "fetch_pending_raws_data",
		"status": "in_progress",
	}).Info("Fetching pending raws data")

	if err := r.db.WithContext(timeoutCtx).
		Joins("JOIN alats ON alats.id = raws.alat_id").
		Select("raws.*, alats.nama").
		Where("raws.is_validated = ?", false).
		Limit(limit).
		Offset(offset).
		Find(&raws).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "fetch_pending_raws_data",
			"status": "error",
			"error":  err.Error(),
		}).Error("Failed to get pending raws data")
		return nil, err
	}
	return raws, nil
}

func (r *validationRepository) UpdateRawStatus(ctx context.Context, rawID uint) error {
	result := r.db.WithContext(ctx).Model(&model.Raw{}).
		Where("id = ?", rawID).
		Update("is_validated", true)

	r.logger.WithFields(logrus.Fields{
		"action": "update_raw_status",
		"raw_id": rawID,
		"status": "updated",
	}).Info("Raw data status updated")

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "update_raw_status",
			"raw_id": rawID,
			"status": "error",
			"error":  result.Error.Error(),
		}).Error("Error updating raw data status")
		return result.Error
	}
	return nil
}

func (r *validationRepository) InsertValidData(ctx context.Context, valid model.Valid) error {
	result := r.db.WithContext(ctx).Create(&valid)
	r.logger.WithFields(logrus.Fields{
		"action": "insert_valid_data",
		"raw_id": valid.RawID,
		"status": "inserting",
	}).Info("Inserting validation data")

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "insert_valid_data",
			"raw_id": valid.RawID,
			"status": "error",
			"error":  result.Error.Error(),
		}).Error("Error inserting validation data")
		return result.Error
	}

	r.logger.WithFields(logrus.Fields{
		"action": "insert_valid_data",
		"raw_id": valid.RawID,
		"status": "success",
	}).Info("Validation data inserted successfully")
	return nil
}

func (r *validationRepository) ValidateAndInsert(ctx context.Context, valid model.Valid) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&valid).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Raw{}).
			Where("id = ?", valid.RawID).
			Update("is_validated", true).
			Error; err != nil {
			return err
		}

		return nil
	})
}
