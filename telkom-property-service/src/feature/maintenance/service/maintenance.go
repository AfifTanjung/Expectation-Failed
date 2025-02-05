package service

import (
	"context"
	"telkom-property-service/src/feature/maintenance/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MaintenanceService interface {
	Create(ctx context.Context, maintenance *model.Maintenance) error
	Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error)
	GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error)
}

type maintenanceService struct {
	repo   repository.MaintenanceRepository
	logger *logrus.Logger
}

func NewMaintenanceService(repo repository.MaintenanceRepository, logger *logrus.Logger) MaintenanceService {
	return &maintenanceService{repo: repo, logger: logger}
}

func (m *maintenanceService) Create(ctx context.Context, maintenance *model.Maintenance) error {
	m.logger.WithFields(logrus.Fields{
		"ID":      maintenance.ID,
		"Start":   maintenance.Start,
		"End":     maintenance.End,
		"Nama":    maintenance.Nama,
		"Laporan": maintenance.Laporan,
		"Status":  maintenance.Status,
	}).Info("[CREATE] Creating Maintenance Log")

	err := m.repo.Create(ctx, maintenance)
	if err != nil {
		m.logger.WithFields(logrus.Fields{
			"ID":      maintenance.ID,
			"Start":   maintenance.Start,
			"End":     maintenance.End,
			"Nama":    maintenance.Nama,
			"Laporan": maintenance.Laporan,
			"Status":  maintenance.Status,
		}).WithError(err).Error("[CREATE] Error creating Maintenance Log")
		return err
	}

	m.logger.WithFields(logrus.Fields{
		"ID":      maintenance.ID,
		"Start":   maintenance.Start,
		"End":     maintenance.End,
		"Nama":    maintenance.Nama,
		"Laporan": maintenance.Laporan,
		"Status":  maintenance.Status,
	}).Info("[CREATE] Maintenance Log Successfully Created")
	return nil
}

func (m *maintenanceService) Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error {
	m.logger.WithField("MaintenanceID", id).Info("[UPDATE] Updating Maintenance Log")

	err := m.repo.Update(ctx, id, maintenance)
	if err != nil {
		m.logger.WithField("MaintenanceID", id).WithError(err).Error("[UPDATE] Error updating Maintenance Log")
		return err
	}

	m.logger.WithField("MaintenanceID", id).Info("[UPDATE] Maintenance Log Successfully Updated")
	return nil
}

func (m *maintenanceService) Delete(ctx context.Context, id uuid.UUID) error {
	m.logger.WithField("MaintenanceID", id).Info("[DELETE] Deleting Maintenance Log")

	err := m.repo.Delete(ctx, id)
	if err != nil {
		m.logger.WithField("MaintenanceID", id).WithError(err).Error("[DELETE] Error deleting Maintenance Log")
		return err
	}

	m.logger.WithField("MaintenanceID", id).Info("[DELETE] Maintenance Log Successfully Deleted")
	return nil
}

func (m *maintenanceService) GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error) {
	m.logger.WithField("MaintenanceID", id).Info("[GET] Fetching Maintenance Log")

	maintenance, err := m.repo.GetLogById(ctx, id)
	if err != nil {
		m.logger.WithField("MaintenanceID", id).WithError(err).Error("[GET] Error fetching Maintenance Log")
		return nil, err
	}

	m.logger.WithField("MaintenanceID", id).Info("[GET] Successfully fetched Maintenance Log")
	return maintenance, nil
}

func (m *maintenanceService) GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	m.logger.WithField("Lantai", lantai).Info("[GET] Fetching Maintenance Logs for lantai")

	alats, err := m.repo.GetLogByLantai(ctx, lantai)
	if err != nil {
		m.logger.WithField("Lantai", lantai).WithError(err).Error("[GET] Error fetching Maintenance Logs")
		return nil, err
	}

	m.logger.WithField("Lantai", lantai).Info("[GET] Successfully fetched Maintenance Logs")
	return alats, nil
}
