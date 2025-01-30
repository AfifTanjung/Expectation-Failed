package service

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
	"telkom-property-service/src/feature/valid/repository"
	"telkom-property-service/src/model"
)

type ValidationServices interface {
	ValidateData(ctx context.Context) error
}

type validationService struct {
	repo   repository.ValidationRepository
	logger *log.Logger
}

func NewValidationService(repo repository.ValidationRepository) *validationService {
	logger := log.New(os.Stdout, "ValidationService: ", log.LstdFlags|log.Lshortfile)
	return &validationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *validationService) ValidateData(ctx context.Context) error {
	s.logger.Println("Starting data validation process...")
	raws, err := s.repo.GetPendingRawData(ctx)
	if err != nil {
		s.logger.Printf("Error fetching raw data: %v\n", err)
		return err
	}

	var wg sync.WaitGroup
	ch := make(chan error, len(raws))

	for _, raw := range raws {
		wg.Add(1)
		go func(raw model.Raw) {
			defer wg.Done()
			if err := s.processRaw(ctx, raw); err != nil {
				ch <- err
			}
		}(raw)
	}

	wg.Wait()
	close(ch)

	for err := range ch {
		if err != nil {
			s.logger.Printf("Validation process encountered an error: %v\n", err)
			return err
		}
	}
	s.logger.Println("Data validation process completed successfully.")
	return nil
}

func (s *validationService) processRaw(ctx context.Context, raw model.Raw) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		isValid, valid, reason := s.validateByAlatType(raw)
		if isValid {
			s.logger.Printf("Raw ID %d is valid. Inserting into valid data...\n", raw.ID)
			if err := s.repo.InsertValidData(ctx, valid); err != nil {
				return err
			}
			raw.ValidationStatus = "valid"
		} else {
			s.logger.Printf("Raw ID %d is invalid: %s\n", raw.ID, reason)
			raw.ValidationStatus = "invalid"
		}
		raw.IsValidated = true
		if err := s.repo.UpdateRawStatus(ctx, raw); err != nil {
			return err
		}
		return nil
	}
}

func (s *validationService) validateByAlatType(raw model.Raw) (bool, model.Valid, string) {
	s.logger.Printf("Validating Raw ID %d with Alat Type: %s\n", raw.ID, raw.Alat.Nama)

	switch raw.Alat.Nama {
	case "IAQMS":
		return s.validateIAQMS(raw)
	case "OAQMS":
		return s.validateOAQMS(raw)
	default:
		return false, model.Valid{}, "Unknown alat type"
	}
}

func (s *validationService) validateIAQMS(raw model.Raw) (bool, model.Valid, string) {
	co2, err := parseAndValidate(raw.Co2, 350, 5050)
	if err != nil {
		return false, model.Valid{}, "CO2 out of range for IAQMS"
	}
	pm2_5, err := parseAndValidate(raw.Pm2_5, 0, 250)
	if err != nil {
		return false, model.Valid{}, "PM2.5 out of range for IAQMS"
	}
	temp, err := parseAndValidate(raw.Temp, 0, 80)
	if err != nil {
		return false, model.Valid{}, "Temperature out of range for IAQMS"
	}
	hum, err := parseAndValidate(raw.Hum, 0, 100)
	if err != nil {
		return false, model.Valid{}, "Humidity out of range for IAQMS"
	}

	return true, model.Valid{
		RawID:      raw.ID,
		Co2:        co2,
		Pm2_5:      pm2_5,
		Temp:       temp,
		Hum:        hum,
		WaktuMasuk: raw.WaktuMasuk,
	}, ""
}

func (s *validationService) validateOAQMS(raw model.Raw) (bool, model.Valid, string) {
	co2, err := parseAndValidate(raw.Co2, 350, 5050)
	if err != nil {
		return false, model.Valid{}, "CO2 out of range for OAQMS"
	}
	pm2_5, err := parseAndValidate(raw.Pm2_5, 0, 250)
	if err != nil {
		return false, model.Valid{}, "PM2.5 out of range for OAQMS"
	}
	temp, err := parseAndValidate(raw.Temp, 0, 80)
	if err != nil {
		return false, model.Valid{}, "Temperature out of range for OAQMS"
	}
	hum, err := parseAndValidate(raw.Hum, 0, 100)
	if err != nil {
		return false, model.Valid{}, "Humidity out of range for OAQMS"
	}
	ws, err := parseAndValidate(raw.Ws, 0, 30)
	if err != nil {
		return false, model.Valid{}, "Wind Speed out of range for OAQMS"
	}
	wd, err := parseAndValidate(raw.Wd, 0, 359)
	if err != nil {
		return false, model.Valid{}, "Wind Direction out of range for OAQMS"
	}

	return true, model.Valid{
		RawID:      raw.ID,
		Co2:        co2,
		Pm2_5:      pm2_5,
		Temp:       temp,
		Hum:        hum,
		Ws:         ws,
		Wd:         wd,
		WaktuMasuk: raw.WaktuMasuk,
	}, ""
}

func parseAndValidate(value string, min, max float64) (string, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil || v < min || v > max {
		return "", errors.New("value out of range")
	}
	return value, nil
}
