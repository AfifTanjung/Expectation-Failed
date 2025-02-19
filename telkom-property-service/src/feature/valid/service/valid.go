package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"telkom-property-service/src/feature/valid/repository"
	"telkom-property-service/src/model"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrValidation    = errors.New("validation error")
	ErrInvalidRange  = errors.New("value out of range")
	ErrInvalidFormat = errors.New("invalid format")
)

type ValidationServices interface {
	ValidateData(ctx context.Context) error
	ValidateBatch(ctx context.Context, batchSize int) error
}

type validationService struct {
	repo    repository.ValidationRepository
	logger  *logrus.Logger
	workers int
}

type SensorRange struct {
	Min, Max float64
	Name     string
}

var (
	iaqmsRanges = []SensorRange{
		{350, 5050, "CO2"},
		{0, 250, "PM2.5"},
		{0, 80, "Temperature"},
		{0, 100, "Humidity"},
	}

	oaqmsRanges = []SensorRange{
		{350, 5050, "CO2"},
		{0, 250, "PM2.5"},
		{0, 80, "Temperature"},
		{0, 100, "Humidity"},
		{0, 30, "Wind Speed"},
		{0, 359, "Wind Direction"},
	}
)

func NewValidationService(repo repository.ValidationRepository, logger *logrus.Logger) ValidationServices {
	return &validationService{
		repo:    repo,
		logger:  logger,
		workers: 5,
	}
}

func (s *validationService) ValidateData(ctx context.Context) error {
	return s.ValidateBatch(ctx, 100)
}

func (s *validationService) ValidateBatch(ctx context.Context, batchSize int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.logger.Info("Starting batch validation process...")

	types := []string{model.IAQMS, model.OAQMS}
	errChan := make(chan error, len(types))
	var wg sync.WaitGroup

	for _, alatType := range types {
		wg.Add(1)
		go func(typ string) {
			defer wg.Done()
			if err := s.processAlatType(timeoutCtx, typ, batchSize); err != nil {
				errChan <- fmt.Errorf("error processing %s: %w", typ, err)
			}
		}(alatType)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %v", errs)
	}

	s.logger.Info("Batch validation process completed successfully")
	return nil
}

func (s *validationService) processAlatType(ctx context.Context, alatType string, batchSize int) error {
	offset := 0
	for {
		raws, err := s.repo.GetPendingRawData(ctx, alatType, batchSize, offset)
		if err != nil {
			return fmt.Errorf("failed to fetch raw data: %w", err)
		}

		if len(raws) == 0 {
			break
		}

		if err := s.processBatch(ctx, raws); err != nil {
			return fmt.Errorf("failed to process batch: %w", err)
		}

		offset += batchSize
	}
	return nil
}

func (s *validationService) processBatch(ctx context.Context, raws []model.Raw) error {
	jobs := make(chan model.Raw, len(raws))
	results := make(chan error, len(raws))
	var wg sync.WaitGroup

	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, jobs, results)
	}

	for _, raw := range raws {
		jobs <- raw
	}
	close(jobs)

	wg.Wait()
	close(results)

	for err := range results {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *validationService) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan model.Raw, results chan<- error) {
	defer wg.Done()
	for raw := range jobs {
		if err := s.processRaw(ctx, raw); err != nil {
			results <- err
			return
		}
	}
}

func (s *validationService) processRaw(ctx context.Context, raw model.Raw) error {
	s.logger.WithFields(logrus.Fields{
		"raw_id":    raw.ID,
		"alat_type": raw.Alat.Type,
	}).Debug("Processing raw data")

	isValid, valid, reason := s.validateByAlatType(raw)
	if !isValid {
		s.logger.WithFields(logrus.Fields{
			"raw_id":    raw.ID,
			"alat_type": raw.Alat.Type,
			"reason":    reason,
			"timestamp": time.Now(),
		}).Warn("Invalid data detected")
		return nil
	}

	if err := s.repo.ValidateAndInsert(ctx, valid); err != nil {
		s.logger.WithFields(logrus.Fields{
			"raw_id": raw.ID,
			"error":  err,
		}).Error("Failed to validate and insert")
		return fmt.Errorf("failed to validate and insert: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"raw_id": raw.ID,
	}).Info("Successfully processed raw data")

	return nil
}

func (s *validationService) validateByAlatType(raw model.Raw) (bool, model.Valid, string) {
	s.logger.WithFields(logrus.Fields{
		"raw_id":    raw.ID,
		"alat_id":   raw.AlatID,
		"alat_type": raw.Alat.Type,
	}).Debug("Validating raw data")

	if raw.Alat.ID == uuid.Nil {
		return false, model.Valid{}, "Alat data not loaded properly"
	}

	if raw.Alat.Type == "" {
		return false, model.Valid{}, fmt.Sprintf("Missing alat type for raw_id: %d", raw.ID)
	}

	switch raw.Alat.Type {
	case model.IAQMS:
		return s.validateIAQMS(raw)
	case model.OAQMS:
		return s.validateOAQMS(raw)
	default:
		return false, model.Valid{}, fmt.Sprintf("Unsupported alat type: %s", raw.Alat.Type)
	}
}

func (s *validationService) validateIAQMS(raw model.Raw) (bool, model.Valid, string) {
	return validateWithRanges(raw, iaqmsRanges)
}

func (s *validationService) validateOAQMS(raw model.Raw) (bool, model.Valid, string) {
	return validateWithRanges(raw, oaqmsRanges)
}

func validateWithRanges(raw model.Raw, ranges []SensorRange) (bool, model.Valid, string) {
	validData := model.Valid{
		RawID:      raw.ID,
		WaktuMasuk: raw.WaktuMasuk,
	}

	values := map[string]string{
		"CO2":            raw.Co2,
		"PM2.5":          raw.Pm2_5,
		"Temperature":    raw.Temp,
		"Humidity":       raw.Hum,
		"Wind Speed":     raw.Ws,
		"Wind Direction": raw.Wd,
	}

	for _, r := range ranges {
		value, exists := values[r.Name]
		if !exists {
			continue
		}

		parsed, err := parseAndValidateRange(value, r)
		if err != nil {
			return false, model.Valid{}, fmt.Sprintf("%s: %v", r.Name, err)
		}

		switch r.Name {
		case "CO2":
			validData.Co2 = parsed
		case "PM2.5":
			validData.Pm2_5 = parsed
		case "Temperature":
			validData.Temp = parsed
		case "Humidity":
			validData.Hum = parsed
		case "Wind Speed":
			validData.Ws = parsed
		case "Wind Direction":
			validData.Wd = parsed
		}
	}

	return true, validData, ""
}

func parseAndValidateRange(value string, r SensorRange) (string, error) {
	if value == "" {
		return "", fmt.Errorf("%w: empty value for %s", ErrInvalidFormat, r.Name)
	}

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", fmt.Errorf("%w: invalid number format for %s: %v", ErrInvalidFormat, r.Name, err)
	}

	if v < r.Min || v > r.Max {
		return "", fmt.Errorf("%w: %s value %v outside range [%v, %v]",
			ErrInvalidRange, r.Name, v, r.Min, r.Max)
	}

	return value, nil
}
