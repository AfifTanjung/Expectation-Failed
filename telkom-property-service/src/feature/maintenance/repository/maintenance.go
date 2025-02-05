package repository

import (
	"context"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MaintenanceRepository interface {
	Create(ctx context.Context, maintenance *model.Maintenance) error
	Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error)
	GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error)
}

type maintenanceRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewMaintenanceRepository(db *gorm.DB, logger *logrus.Logger) MaintenanceRepository {
	return &maintenanceRepository{db: db, logger: logger}
}

func (r *maintenanceRepository) Create(ctx context.Context, maintenance *model.Maintenance) error {
	r.logger.WithField("AlatID", maintenance.AlatID).Info("[Create] Inserting maintenance record")
	if err := r.db.WithContext(ctx).Create(maintenance).Error; err != nil {
		r.logger.WithField("AlatID", maintenance.AlatID).WithError(err).Error("[Create] Error inserting maintenance")
		return err
	}
	r.logger.WithField("AlatID", maintenance.AlatID).Info("[Create] Maintenance record inserted successfully")
	return nil
}

func (r *maintenanceRepository) Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error {
	r.logger.WithField("MaintenanceID", id).Info("[Update] Updating maintenance record")
	if err := r.db.WithContext(ctx).Model(&model.Maintenance{}).
		Where("id = ?", id).
		Updates(maintenance).Error; err != nil {
		r.logger.WithField("MaintenanceID", id).WithError(err).Error("[Update] Error updating maintenance")
		return err
	}
	r.logger.WithField("MaintenanceID", id).Info("[Update] Maintenance record updated successfully")
	return nil
}

func (r *maintenanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.WithField("MaintenanceID", id).Info("[Delete] Deleting maintenance record")
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Maintenance{}).Error; err != nil {
		r.logger.WithField("MaintenanceID", id).WithError(err).Error("[Delete] Error deleting maintenance")
		return err
	}
	r.logger.WithField("MaintenanceID", id).Info("[Delete] Maintenance record deleted successfully")
	return nil
}

func (r *maintenanceRepository) GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.WithField("Lantai", lantai).Info("[GetLogByLantai] Fetching maintenance logs")
	err := r.db.WithContext(ctx).
		Joins("JOIN alats ON alats.id = maintenances.alat_id").
		Where("alats.lantai = ?", lantai).
		Find(&alats).Error
	if err != nil {
		r.logger.WithField("Lantai", lantai).WithError(err).Error("[GetLogByLantai] Error fetching logs")
		return nil, err
	}
	r.logger.WithField("Lantai", lantai).WithField("TotalLogs", len(alats)).Info("[GetLogByLantai] Found maintenance logs")
	return alats, nil
}

func (r *maintenanceRepository) GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error) {
	var maintenance model.Maintenance
	r.logger.WithField("MaintenanceID", id).Info("[GetLogById] Fetching maintenance log")
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&maintenance).Error
	if err != nil {
		r.logger.WithField("MaintenanceID", id).WithError(err).Error("[GetLogById] Error fetching log")
		return nil, err
	}
	r.logger.WithField("MaintenanceID", id).Info("[GetLogById] Maintenance log fetched successfully")
	return &maintenance, nil
}
