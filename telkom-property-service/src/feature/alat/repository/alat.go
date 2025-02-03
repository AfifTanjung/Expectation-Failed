package repository

import (
	"context"
	"log"
	"telkom-property-service/src/model"

	"gorm.io/gorm"
)

type AlatRepository interface {
	Create(ctx context.Context, alat model.Alat) error
	GetByLantai(ctx context.Context, lantai int) ([]model.Alat, error)
	GetByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error)
}

type alatRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewAlatRepository(db *gorm.DB, logger *log.Logger) AlatRepository {
	return &alatRepository{db: db, logger: logger}
}

func (r *alatRepository) Create(ctx context.Context, alat model.Alat) error {
	r.logger.Printf("[Create] Creating alat: ID=%s, Nama=%s", alat.ID, alat.Nama)
	r.logger.Printf("Data before create: %+v\n", alat)
	if err := r.db.WithContext(ctx).Create(&alat).Error; err != nil {
		r.logger.Printf("[Create] Error creating alat: %v", err)
		return err
	}
	r.logger.Printf("[Create] Alat created successfully: ID=%s, Nama=%s", alat.ID, alat.Nama)
	return nil
}

func (r *alatRepository) GetByLantai(ctx context.Context, lantai int) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.Printf("[GetByLantai] Fetching alat for lantai=%d", lantai)
	if err := r.db.WithContext(ctx).Preload("Raws").Where("lantai = ?", lantai).Find(&alats).Error; err != nil {
		r.logger.Printf("[GetByLantai] Error fetching alat: %v", err)
		return nil, err
	}
	for _, alat := range alats {
		r.logger.Printf("[GetByLantai] Alat ID=%v has %d raws", alat.ID, len(alat.Raws))
	}

	r.logger.Printf("[GetByLantai] Successfully fetched %d alat(s)", len(alats))
	return alats, nil
}

func (r *alatRepository) GetByStasiun(ctx context.Context, stasiun string) ([]model.Alat, error) {
	var alats []model.Alat
	r.logger.Printf("[GetByStasiun] Fetching alat for stasiun=%s", stasiun)
	if err := r.db.WithContext(ctx).Where("stasiun = ?", stasiun).Find(&alats).Error; err != nil {
		r.logger.Printf("[GetByStasiun] Error fetching alat: %v", err)
		return nil, err
	}
	r.logger.Printf("[GetByStasiun] Successfully fetched %d alat(s)", len(alats))
	return alats, nil
}
