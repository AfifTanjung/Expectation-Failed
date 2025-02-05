package repository

import (
	"context"
	"telkom-property-service/src/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AlatRepository interface {
	Create(ctx context.Context, alat model.Alat) error
	GetByLantai(ctx context.Context, lantai int) ([]model.Alat, error)
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
	r.logger.WithFields(logrus.Fields{
		"ID":   alat.ID,
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
		"Nama": alat.Nama,
	}).Info("[Create] Alat created successfully")
	return nil
}

func (r *alatRepository) GetByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.WithField("lantai", lantai).Info("[GetByLantai] Fetching alat")

	if err := r.db.WithContext(ctx).Preload("Raws").Where("lantai = ?", lantai).Find(&alats).Error; err != nil {
		r.logger.WithError(err).Error("[GetByLantai] Error fetching alat")
		return nil, err
	}

	for _, alat := range alats {
		r.logger.WithFields(logrus.Fields{
			"ID":   alat.ID,
			"Raws": len(alat.Raws),
		}).Debug("[GetByLantai] Alat fetched")
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
