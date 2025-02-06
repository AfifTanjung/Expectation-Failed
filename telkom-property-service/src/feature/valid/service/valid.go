package service

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"telkom-property-service/src/feature/valid/repository"
	"telkom-property-service/src/model"
	"time"

	"github.com/sirupsen/logrus"
)

type ValidationServices interface {
	ValidateData(ctx context.Context) error
}

type validationService struct {
	repo   repository.ValidationRepository
	logger *logrus.Logger
}

func NewValidationService(repo repository.ValidationRepository, logger *logrus.Logger) *validationService {
	return &validationService{
		repo:   repo,
		logger: logger,
	}
}

func (s *validationService) ValidateData(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	s.logger.Info("Starting data validation process...")

	// Ambil data berdasarkan tipe alat (IAQMS dan OAQMS)
	types := []string{model.IAQMS, model.OAQMS} // Gunakan model.AlatType
	var allRaws []model.Raw
	for _, alatType := range types {
		raws, err := s.repo.GetPendingRawData(timeoutCtx, alatType, 100, 0)
		if err != nil {
			s.logger.WithError(err).Errorf("Error fetching raw data for type: %s", alatType)
			return err
		}
		allRaws = append(allRaws, raws...)
	}

	var wg sync.WaitGroup
	ch := make(chan error, len(allRaws))
	for _, raw := range allRaws {
		wg.Add(1)
		go func(raw model.Raw) {
			defer wg.Done()
			if err := s.processRaw(timeoutCtx, raw); err != nil {
				ch <- err
			}
		}(raw)
	}
	wg.Wait()
	close(ch)

	for err := range ch {
		if err != nil {
			s.logger.WithError(err).Error("Validation process encountered an error")
			return err
		}
	}
	s.logger.Info("Data validation process completed successfully.")
	return nil
}

func (s *validationService) processRaw(ctx context.Context, raw model.Raw) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		isValid, valid, reason := s.validateByAlatType(raw)
		if isValid {
			s.logger.Infof("Raw ID %d is valid. Inserting into valid data...", raw.ID)
			if err := s.repo.InsertValidData(ctx, valid); err != nil {
				return err
			}
		} else {
			s.logger.Warnf("Raw ID %d is invalid: %s", raw.ID, reason)
		}
		return nil
	}
}

func (s *validationService) validateByAlatType(raw model.Raw) (bool, model.Valid, string) {
	s.logger.Infof("Validating Raw ID %d with Alat Type: %s", raw.ID, raw.Alat.Type)
	switch raw.Alat.Type {
	case model.IAQMS:
		return s.validateIAQMS(raw)
	case model.OAQMS:
		return s.validateOAQMS(raw)
	default:
		return false, model.Valid{}, "Unknown alat Type"
	}
}

func (s *validationService) validateIAQMS(raw model.Raw) (bool, model.Valid, string) {
	return validateSensorData(raw, 350, 5050, 0, 250, 0, 80, 0, 100)
}

func (s *validationService) validateOAQMS(raw model.Raw) (bool, model.Valid, string) {
	return validateSensorData(raw, 350, 5050, 0, 250, 0, 80, 0, 100, 0, 30, 0, 359)
}

func validateSensorData(raw model.Raw, ranges ...float64) (bool, model.Valid, string) {
	values := []string{raw.Co2, raw.Pm2_5, raw.Temp, raw.Hum, raw.Ws, raw.Wd}
	validData := model.Valid{RawID: raw.ID, WaktuMasuk: raw.WaktuMasuk}
	for i, val := range values {
		if i >= len(ranges)/2 {
			break
		}
		parsed, err := parseAndValidate(val, ranges[i*2], ranges[i*2+1])
		if err != nil {
			return false, model.Valid{}, err.Error()
		}
		setField(&validData, i, parsed)
	}
	return true, validData, ""
}

func parseAndValidate(value string, min, max float64) (string, error) {
	v, err := strconv.ParseFloat(value, 64)
	if err != nil || v < min || v > max {
		return "", errors.New("value out of range")
	}
	return value, nil
}

func setField(validData *model.Valid, index int, value string) {
	switch index {
	case 0:
		validData.Co2 = value
	case 1:
		validData.Pm2_5 = value
	case 2:
		validData.Temp = value
	case 3:
		validData.Hum = value
	case 4:
		validData.Ws = value
	case 5:
		validData.Wd = value
	}
}
