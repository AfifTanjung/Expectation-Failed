package services

import (
	"context"
	"errors"
	"log"
	"telkom-property-service/src/feature/alat/repository"
	"telkom-property-service/src/model"

	"github.com/google/uuid"
)

type AlatService interface {
	CreateAlat(ctx context.Context, alat model.Alat) error
	GetAlatByLantai(ctx context.Context, lantai int) ([]model.Alat, error)
	GetAlatByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error)
}

type alatService struct {
	repo   repository.AlatRepository
	logger *log.Logger
}

func NewAlatService(repo repository.AlatRepository, logger *log.Logger) AlatService {
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
	s.logger.Printf("[CreateAlat] Starting creation of alat: Nama=%s, Lantai=%d, Stasiun=%s", alat.Nama, alat.Lantai, alat.Stasiun)
	if alat.ID == uuid.Nil {
		alat.ID = uuid.New()
	}
	err := validateAlat(alat)
	if err != nil {
		s.logger.Println("[CreateAlat] Validation failed:", err)
		return err
	}
	err = s.repo.Create(ctx, alat)
	if err != nil {
		s.logger.Println("[CreateAlat] Error creating alat:", err)
		return err
	}
	s.logger.Printf("[CreateAlat] Successfully created alat: ID=%s, Nama=%s", alat.ID, alat.Nama)
	return nil
}

func (s *alatService) GetAlatByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	s.logger.Printf("[GetAlatByLantai] Starting fetch for lantai=%d", lantai)
	if lantai <= 0 {
		err := errors.New("lantai harus lebih dari 0")
		s.logger.Println("[GetAlatByLantai] Validation failed:", err)
		return nil, err
	}
	alats, err := s.repo.GetByLantai(ctx, lantai)
	if err != nil {
		s.logger.Println("[GetAlatByLantai] Error fetching alat by lantai:", err)
		return nil, err
	}
	if len(alats) == 0 {
		s.logger.Printf("[GetAlatByLantai] No alat found for lantai=%d", lantai)
	} else {
		s.logger.Printf("[GetAlatByLantai] Found %d alat(s) for lantai=%d", len(alats), lantai)
	}
	return alats, nil
}

func (s *alatService) GetAlatByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error) {
	s.logger.Printf("[GetAlatByStasiun] Starting fetch for stasiun=%s", stasiun)

	if stasiun == "" {
		err := errors.New("stasiun tidak boleh kosong")
		s.logger.Println("[GetAlatByStasiun] Validation failed:", err)
		return nil, err
	}
	alats, err := s.repo.GetByStasiun(ctx, stasiun)
	if err != nil {
		s.logger.Println("[GetAlatByStasiun] Error fetching alat by stasiun:", err)
		return nil, err
	}
	if len(alats) == 0 {
		s.logger.Printf("[GetAlatByStasiun] No alat found for stasiun=%s", stasiun)
	} else {
		s.logger.Printf("[GetAlatByStasiun] Found %d alat(s) for stasiun=%s", len(alats), stasiun)
	}

	return alats, nil
}
