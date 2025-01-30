package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"telkom-property-service/src/feature/raw/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
)

var (
	ErrEmptyAlatID   = errors.New("alatID tidak boleh kosong")
	ErrAlatNotFound  = errors.New("alat tidak ditemukan")
	ErrFailedToFetch = errors.New("gagal mengambil data")
)

type RawServices interface {
	CreateByAlatId(ctx context.Context, raw model.Raw) error
	GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error)
}

type rawService struct {
	repo   repository.RawRepository
	logger *log.Logger
}

func NewRawService(repo repository.RawRepository, logger *log.Logger) RawServices {
	return &rawService{repo: repo, logger: logger}
}

func (s *rawService) CreateByAlatId(ctx context.Context, raw model.Raw) error {
	s.logger.Printf("[INFO] [CreateByAlatId] Starting processing for AlatID: %s", raw.AlatID)
	if raw.AlatID == uuid.Nil {
		s.logger.Println("[ERROR] [CreateByAlatId] Validation failed: ", ErrEmptyAlatID)
		return fmt.Errorf("validation error: %w", ErrEmptyAlatID)
	}
	if err := s.repo.CreateByAlatId(ctx, raw); err != nil {
		s.logger.Printf("[ERROR] [CreateByAlatId] Failed to create Raw data for AlatID: %s, Error: %v", raw.AlatID, err)
		return fmt.Errorf("failed to create raw data: %w", err)
	}

	s.logger.Printf("[INFO] [CreateByAlatId] Successfully created Raw data for AlatID: %s", raw.AlatID)
	return nil
}

func (s *rawService) GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error) {
	s.logger.Printf("[INFO] [GetById] Fetching data for AlatID: %s", alatID)
	if alatID == uuid.Nil {
		s.logger.Println("[ERROR] [GetById] Validation failed: ", ErrEmptyAlatID)
		return nil, nil, fmt.Errorf("validation error: %w", ErrEmptyAlatID)
	}
	alat, raws, err := s.repo.GetById(ctx, alatID)
	if err != nil {
		s.logger.Printf("[ERROR] [GetById] Failed to fetch data for AlatID: %s, Error: %v", alatID, err)
		return nil, nil, fmt.Errorf("failed to fetch data for AlatID %s: %w", alatID, err)
	}
	if alat == nil {
		s.logger.Printf("[ERROR] [GetById] Alat not found for AlatID: %s", alatID)
		return nil, nil, fmt.Errorf("alat with ID %s not found: %w", alatID, ErrAlatNotFound)
	}
	s.logger.Printf("[INFO] [GetById] Successfully fetched data for AlatID: %s, Total records: %d", alatID, len(raws))
	return alat, raws, nil
}
