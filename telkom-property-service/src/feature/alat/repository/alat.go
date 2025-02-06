package repository

import (
	"context"
	"fmt"
	"telkom-property-service/src/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AlatRepository interface {
	Create(ctx context.Context, alat model.Alat) error
	GetByLantai(ctx context.Context, lantai int, limit, offset int) ([]model.Alat, error)
	GetByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error)
}

type alatRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAlatRepository(db *gorm.DB, logger *logrus.Logger) AlatRepository {
	return &alatRepository{db: db, logger: logger}
}

func (r *alatRepository) Create(ctx context.Context, alat model.Alat) error {
	if alat.Nama == "" || alat.Lantai == 0 || alat.Stasiun == "" {
		return fmt.Errorf("nama, lantai, and stasiun are required fields")
	}
	r.logger.WithFields(logrus.Fields{
		"ID":   alat.ID,
		"type": alat.Type,
		"Nama": alat.Nama,
	}).Info("[Create] Creating alat")

	r.logger.WithFields(logrus.Fields{
		"alat": alat,
	}).Debug("[Create] Data before create")

	if err := r.db.WithContext(ctx).Create(&alat).Error; err != nil {
		r.logger.WithError(err).Error("[Create] Error creating alat")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"ID":   alat.ID,
		"type": alat.Type,
		"Nama": alat.Nama,
	}).Info("[Create] Alat created successfully")
	return nil
}

func (r *alatRepository) GetByLantai(ctx context.Context, lantai int, limit, offset int) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.WithField("lantai", lantai).Info("[GetByLantai] Fetching alat")

	query := r.db.WithContext(ctx).Preload("Raws").Where("lantai = ?", lantai)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&alats).Error; err != nil {
		r.logger.WithError(err).Error("[GetByLantai] Error fetching alat")
		return nil, fmt.Errorf("failed to fetch alat by lantai: %w", err)
	}

	r.logger.WithField("count", len(alats)).Info("[GetByLantai] Successfully fetched alat(s)")
	return alats, nil
}

func (r *alatRepository) GetByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.WithField("stasiun", stasiun).Info("[GetByStasiun] Fetching alat")

	if err := r.db.WithContext(ctx).Where("stasiun = ?", stasiun).Find(&alats).Error; err != nil {
		r.logger.WithError(err).Error("[GetByStasiun] Error fetching alat")
		return nil, err
	}

	r.logger.WithField("count", len(alats)).Info("[GetByStasiun] Successfully fetched alat(s)")
	return alats, nil
}
