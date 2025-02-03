package repository

import (
	"context"
	"log"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
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
	logger *log.Logger
}

func NewMaintenanceRepository(db *gorm.DB, logger *log.Logger) MaintenanceRepository {
	return &maintenanceRepository{db: db, logger: logger}
}

func (r *maintenanceRepository) Create(ctx context.Context, maintenance *model.Maintenance) error {
	r.logger.Println("[Create] Inserting maintenance record")
	if err := r.db.WithContext(ctx).Create(maintenance).Error; err != nil {
		r.logger.Printf("[Create] Error inserting maintenance: %v", err)
		return err
	}
	r.logger.Println("[Create] Maintenance record inserted successfully")
	return nil
}

func (r *maintenanceRepository) Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error {
	r.logger.Printf("[Update] Updating maintenance record with ID: %v", id)
	if err := r.db.WithContext(ctx).Model(&model.Maintenance{}).
		Where("id = ?", id).
		Updates(maintenance).Error; err != nil {
		r.logger.Printf("[Update] Error updating maintenance: %v", err)
		return err
	}
	r.logger.Println("[Update] Maintenance record updated successfully")
	return nil
}

func (r *maintenanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Printf("[Delete] Deleting maintenance record with ID: %v", id)
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Maintenance{}).Error; err != nil {
		r.logger.Printf("[Delete] Error deleting maintenance: %v", err)
		return err
	}
	r.logger.Println("[Delete] Maintenance record deleted successfully")
	return nil
}

func (r *maintenanceRepository) GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.Printf("[GetLogByLantai] Fetching maintenance logs for lantai: %d", lantai)
	err := r.db.WithContext(ctx).
		Joins("JOIN alats ON alats.id = maintenances.alat_id").
		Where("alats.lantai = ?", lantai).
		Find(&alats).Error
	if err != nil {
		r.logger.Printf("[GetLogByLantai] Error fetching logs: %v", err)
		return nil, err
	}
	r.logger.Printf("[GetLogByLantai] Found %d logs", len(alats))
	return alats, nil
}

func (r *maintenanceRepository) GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error) {
	var maintenance model.Maintenance
	r.logger.Printf("[GetLogById] Fetching maintenance log for ID: %v", id)
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&maintenance).Error
	if err != nil {
		r.logger.Printf("[GetLogById] Error fetching log: %v", err)
		return nil, err
	}
	r.logger.Println("[GetLogById] Maintenance log fetched successfully")
	return &maintenance, nil
}
