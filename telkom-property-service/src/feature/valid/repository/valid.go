package repository

import (
	"context"
	"fmt"
	"telkom-property-service/src/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Custom error types
type ValidationError struct {
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

// Metrics struct untuk monitoring
type Metrics struct {
	QueryDurations map[string][]time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{
		QueryDurations: make(map[string][]time.Duration),
	}
}

func (m *Metrics) ObserveQueryDuration(operation string, duration time.Duration) {
	m.QueryDurations[operation] = append(m.QueryDurations[operation], duration)
}

// Repository interface
type ValidationRepository interface {
	GetPendingRawData(ctx context.Context, alatType string, limit, offset int) ([]model.Raw, error)
	UpdateRawStatus(ctx context.Context, rawID uint) error
	BatchUpdateRawStatus(ctx context.Context, rawIDs []uint) error
	InsertValidData(ctx context.Context, valid model.Valid) error
	ValidateAndInsert(ctx context.Context, valid model.Valid) error
	GetValidationMetrics() map[string][]time.Duration
}

type validationRepository struct {
	db      *gorm.DB
	logger  *logrus.Logger
	metrics *Metrics
}

var _ ValidationRepository = (*validationRepository)(nil)

func NewValidationRepository(db *gorm.DB, logger *logrus.Logger) ValidationRepository {
	return &validationRepository{
		db:      db,
		logger:  logger,
		metrics: NewMetrics(),
	}
}

const (
	defaultQueryTimeout = 5 * time.Second
	maxRetries          = 3
)

func (r *validationRepository) GetPendingRawData(ctx context.Context, alatType string, limit, offset int) ([]model.Raw, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var raws []model.Raw

	r.logger.WithFields(logrus.Fields{
		"alat_type": alatType,
		"limit":     limit,
		"offset":    offset,
	}).Debug("Fetching pending raw data")

	result := r.db.WithContext(timeoutCtx).
		Joins("INNER JOIN alats ON alats.id = raws.alat_id").
		Preload("Alat").
		Select("raws.*, alats.type as alat_type").
		Where("raws.is_validated = ? AND alats.type = ?", false, alatType).
		Limit(limit).
		Offset(offset).
		Find(&raws)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"error":     result.Error,
			"alat_type": alatType,
		}).Error("Failed to fetch pending raw data")
		return nil, result.Error
	}

	for _, raw := range raws {
		r.logger.WithFields(logrus.Fields{
			"raw_id":    raw.ID,
			"alat_id":   raw.AlatID,
			"alat_type": raw.Alat.Type,
		}).Debug("Raw data fetched")
	}

	r.logger.WithFields(logrus.Fields{
		"count":     len(raws),
		"alat_type": alatType,
	}).Info("Successfully fetched pending raw data")

	return raws, nil
}

func (r *validationRepository) UpdateRawStatus(ctx context.Context, rawID uint) error {
	if rawID == 0 {
		return &ValidationError{"invalid raw ID", nil}
	}

	start := time.Now()
	defer func() {
		r.metrics.ObserveQueryDuration("update_raw_status", time.Since(start))
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result := r.db.WithContext(timeoutCtx).Model(&model.Raw{}).
		Where("id = ?", rawID).
		Update("is_validated", true)

	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "update_raw_status",
			"raw_id": rawID,
			"status": "error",
			"error":  result.Error.Error(),
		}).Error("Error updating raw data status")
		return &ValidationError{"failed to update raw status", result.Error}
	}

	if result.RowsAffected == 0 {
		return &ValidationError{"no rows affected", nil}
	}

	r.logger.WithFields(logrus.Fields{
		"action":      "update_raw_status",
		"raw_id":      rawID,
		"status":      "updated",
		"duration_ms": time.Since(start).Milliseconds(),
	}).Info("Raw data status updated")
	return nil
}

func (r *validationRepository) BatchUpdateRawStatus(ctx context.Context, rawIDs []uint) error {
	if len(rawIDs) == 0 {
		return &ValidationError{"empty raw IDs slice", nil}
	}

	start := time.Now()
	defer func() {
		r.metrics.ObserveQueryDuration("batch_update_raw_status", time.Since(start))
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result := r.db.WithContext(timeoutCtx).
		Model(&model.Raw{}).
		Where("id IN ?", rawIDs).
		Update("is_validated", true)

	if result.Error != nil {
		return &ValidationError{"failed to batch update raw status", result.Error}
	}

	r.logger.WithFields(logrus.Fields{
		"action":        "batch_update_raw_status",
		"rows_affected": result.RowsAffected,
		"duration_ms":   time.Since(start).Milliseconds(),
	}).Info("Batch update completed")

	return nil
}

func (r *validationRepository) InsertValidData(ctx context.Context, valid model.Valid) error {
	if valid.RawID == 0 {
		return &ValidationError{"invalid raw ID", nil}
	}

	start := time.Now()
	defer func() {
		r.metrics.ObserveQueryDuration("insert_valid_data", time.Since(start))
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	result := r.db.WithContext(timeoutCtx).Create(&valid)
	if result.Error != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "insert_valid_data",
			"raw_id": valid.RawID,
			"status": "error",
			"error":  result.Error.Error(),
		}).Error("Error inserting validation data")
		return &ValidationError{"failed to insert valid data", result.Error}
	}

	r.logger.WithFields(logrus.Fields{
		"action":      "insert_valid_data",
		"raw_id":      valid.RawID,
		"status":      "success",
		"duration_ms": time.Since(start).Milliseconds(),
	}).Info("Validation data inserted successfully")
	return nil
}

func (r *validationRepository) ValidateAndInsert(ctx context.Context, valid model.Valid) error {
	if valid.RawID == 0 {
		return &ValidationError{"rawID cannot be zero", nil}
	}

	start := time.Now()
	defer func() {
		r.metrics.ObserveQueryDuration("validate_and_insert", time.Since(start))
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()

	err := r.db.WithContext(timeoutCtx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&valid).Error; err != nil {
			return &ValidationError{"failed to insert valid data", err}
		}

		if err := tx.Model(&model.Raw{}).
			Where("id = ?", valid.RawID).
			Update("is_validated", true).
			Error; err != nil {
			return &ValidationError{"failed to update raw status", err}
		}

		return nil
	})

	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"action": "validate_and_insert",
			"raw_id": valid.RawID,
			"status": "error",
			"error":  err.Error(),
		}).Error("Transaction failed")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"action":      "validate_and_insert",
		"raw_id":      valid.RawID,
		"status":      "success",
		"duration_ms": time.Since(start).Milliseconds(),
	}).Info("Transaction completed successfully")

	return nil
}

func (r *validationRepository) GetValidationMetrics() map[string][]time.Duration {
	return r.metrics.QueryDurations
}
