package service

import (
	"context"
	"log"
	"telkom-property-service/src/feature/maintenance/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
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
	logger *log.Logger
}

func NewMaintenanceService(repo repository.MaintenanceRepository, logger *log.Logger) MaintenanceService {
	return &maintenanceService{repo: repo, logger: logger}
}

func (m *maintenanceService) Create(ctx context.Context, maintenance *model.Maintenance) error {
	m.logger.Printf("[CREATE] Creating Maintenance Log : ID:%s, \n Start:%s, \n End:%s \n Nama:%s \n Laporan:%s \n Status : %v",
		maintenance.ID, maintenance.Start, maintenance.End, maintenance.Nama, maintenance.Laporan, maintenance.Status)

	err := m.repo.Create(ctx, maintenance)
	if err != nil {
		m.logger.Println("[CREATE] Error Maintenance Log ", err)
		return err
	}
	m.logger.Printf("[CREATE] Maintenance Log Successfully Created: ID:%s, \n Start:%s, \n End:%s \n Nama:%s \n Laporan:%s \n Status : %v",
		maintenance.ID, maintenance.Start, maintenance.End, maintenance.Nama, maintenance.Laporan, maintenance.Status)
	return nil
}

func (m *maintenanceService) Update(ctx context.Context, id uuid.UUID, maintenance *model.Maintenance) error {
	m.logger.Printf("[UPDATE] Updating Maintenance Log with ID: %s", id)

	err := m.repo.Update(ctx, id, maintenance)
	if err != nil {
		m.logger.Println("[UPDATE] Error updating Maintenance Log:", err)
		return err
	}
	m.logger.Printf("[UPDATE] Maintenance Log Successfully Updated: ID:%s", id)
	return nil
}

func (m *maintenanceService) Delete(ctx context.Context, id uuid.UUID) error {
	m.logger.Printf("[DELETE] Deleting Maintenance Log with ID: %s", id)

	err := m.repo.Delete(ctx, id)
	if err != nil {
		m.logger.Println("[DELETE] Error deleting Maintenance Log:", err)
		return err
	}
	m.logger.Printf("[DELETE] Maintenance Log Successfully Deleted: ID:%s", id)
	return nil
}

func (m *maintenanceService) GetLogById(ctx context.Context, id uuid.UUID) (*model.Maintenance, error) {
	m.logger.Printf("[GET] Fetching Maintenance Log for ID: %s", id)

	maintenance, err := m.repo.GetLogById(ctx, id)
	if err != nil {
		m.logger.Println("[GET] Error fetching Maintenance Log:", err)
		return nil, err
	}
	m.logger.Printf("[GET] Successfully fetched Maintenance Log for ID: %s", id)
	return maintenance, nil
}

func (m *maintenanceService) GetLogByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	m.logger.Printf("[GET] Fetching Maintenance Logs for lantai: %d", lantai)

	alats, err := m.repo.GetLogByLantai(ctx, lantai)
	if err != nil {
		m.logger.Println("[GET] Error fetching Maintenance Logs:", err)
		return nil, err
	}
	m.logger.Printf("[GET] Successfully fetched Maintenance Logs for lantai: %d", lantai)
	return alats, nil
}
