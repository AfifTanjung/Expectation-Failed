package services

import (
	"context"
	"errors"
	"fmt"
	"telkom-property-service/src/feature/alat/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AlatService interface {
	CreateAlat(ctx context.Context, alat model.Alat) error
	GetAlatByLantai(ctx context.Context, lantai, limit, offset int) ([]model.Alat, error)
	GetAlatByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error)
}

type alatService struct {
	repo   repository.AlatRepository
	logger *logrus.Logger
}

func NewAlatService(repo repository.AlatRepository, logger *logrus.Logger) AlatService {
	return &alatService{repo: repo, logger: logger}
}

func validateAlat(alat model.Alat) error {
	if alat.Nama == "" {
		return errors.New("nama alat tidak boleh kosong")
	}
	if alat.Stasiun == "" {
		return errors.New("stasiun tidak boleh kosong")
	}
	if alat.Lantai <= 0 {
		return errors.New("lantai harus lebih dari 0")
	}
	return nil
}

func (s *alatService) CreateAlat(ctx context.Context, alat model.Alat) error {
	s.logger.WithFields(logrus.Fields{
		"Nama":    alat.Nama,
		"Lantai":  alat.Lantai,
		"type":    alat.Type,
		"Stasiun": alat.Stasiun,
	}).Info("[CreateAlat] Starting creation of alat")
	if alat.ID == uuid.Nil {
		alat.ID = uuid.New()
	}
	err := validateAlat(alat)
	if err != nil {
		s.logger.WithError(err).Error("[CreateAlat] Validation failed")
		return fmt.Errorf("validation failed: %w", err)
	}
	err = s.repo.Create(ctx, alat)
	if err != nil {
		s.logger.WithError(err).Error("[CreateAlat] Error creating alat")
		return fmt.Errorf("failed to create alat: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"ID":   alat.ID,
		"Nama": alat.Nama,
	}).Info("[CreateAlat] Successfully created alat")
	return nil
}

func (s *alatService) GetAlatByLantai(ctx context.Context, lantai, limit, offset int) ([]model.Alat, error) {
	s.logger.WithFields(logrus.Fields{
		"lantai": lantai,
		"limit":  limit,
		"offset": offset,
	}).Info("[GetAlatByLantai] Starting fetch for lantai")

	if lantai <= 0 {
		err := errors.New("lantai harus lebih dari 0")
		s.logger.WithError(err).Error("[GetAlatByLantai] Validation failed")
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	alats, err := s.repo.GetByLantai(ctx, lantai, limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("[GetAlatByLantai] Error fetching alat by lantai")
		return nil, fmt.Errorf("failed to fetch alat by lantai: %w", err)
	}

	if len(alats) == 0 {
		s.logger.WithField("lantai", lantai).Info("[GetAlatByLantai] No alat found for lantai")
	} else {
		s.logger.WithFields(logrus.Fields{
			"lantai": lantai,
			"count":  len(alats),
		}).Info("[GetAlatByLantai] Found alat(s) for lantai")
	}

	return alats, nil
}

func (s *alatService) GetAlatByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error) {
	s.logger.WithField("stasiun", stasiun).Info("[GetAlatByStasiun] Starting fetch for stasiun")

	if stasiun == "" {
		err := errors.New("stasiun tidak boleh kosong")
		s.logger.WithError(err).Error("[GetAlatByStasiun] Validation failed")
		return nil, fmt.Errorf("invalid input: %w", err)
	}
	alats, err := s.repo.GetByStasiun(ctx, stasiun)
	if err != nil {
		s.logger.WithError(err).Error("[GetAlatByStasiun] Error fetching alat by stasiun")
		return nil, fmt.Errorf("failed to fetch alat by stasiun: %w", err)
	}
	if len(alats) == 0 {
		s.logger.WithField("stasiun", stasiun).Info("[GetAlatByStasiun] No alat found for stasiun")
	} else {
		s.logger.WithFields(logrus.Fields{
			"stasiun": stasiun,
			"count":   len(alats),
		}).Info("[GetAlatByStasiun] Found alat(s) for stasiun")
	}

	return alats, nil
}
