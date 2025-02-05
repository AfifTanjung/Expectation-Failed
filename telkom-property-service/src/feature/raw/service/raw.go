package service

import (
	"context"
	"errors"
	"fmt"
	"telkom-property-service/src/feature/raw/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logger *logrus.Logger
}

func NewRawService(repo repository.RawRepository, logger *logrus.Logger) RawServices {
	return &rawService{repo: repo, logger: logger}
}

func (s *rawService) CreateByAlatId(ctx context.Context, raw model.Raw) error {
	s.logger.WithField("AlatID", raw.AlatID).Info("[CreateByAlatId] Starting processing")
	if raw.AlatID == uuid.Nil {
		s.logger.WithField("AlatID", raw.AlatID).Error("[CreateByAlatId] Validation failed: AlatID cannot be empty")
		return fmt.Errorf("validation error: %w", ErrEmptyAlatID)
	}
	if err := s.repo.CreateByAlatId(ctx, raw); err != nil {
		s.logger.WithField("AlatID", raw.AlatID).WithError(err).Error("[CreateByAlatId] Failed to create raw data")
		return fmt.Errorf("failed to create raw data: %w", err)
	}

	s.logger.WithField("AlatID", raw.AlatID).Info("[CreateByAlatId] Successfully created raw data")
	return nil
}

func (s *rawService) GetById(ctx context.Context, alatID uuid.UUID) (*model.Alat, []model.Raw, error) {
	s.logger.WithField("AlatID", alatID).Info("[GetById] Fetching data")

	if alatID == uuid.Nil {
		s.logger.WithField("AlatID", alatID).Error("[GetById] Validation failed: AlatID cannot be empty")
		return nil, nil, fmt.Errorf("validation error: %w", ErrEmptyAlatID)
	}

	alat, raws, err := s.repo.GetById(ctx, alatID)
	if err != nil {
		s.logger.WithField("AlatID", alatID).WithError(err).Error("[GetById] Failed to fetch data")
		return nil, nil, fmt.Errorf("failed to fetch data for AlatID %s: %w", alatID, err)
	}
	if alat == nil {
		s.logger.WithField("AlatID", alatID).Error("[GetById] Alat not found")
		return nil, nil, fmt.Errorf("alat with ID %s not found: %w", alatID, ErrAlatNotFound)
	}

	s.logger.WithField("AlatID", alatID).WithField("TotalRecords", len(raws)).Info("[GetById] Successfully fetched data")
	return alat, raws, nil
}
