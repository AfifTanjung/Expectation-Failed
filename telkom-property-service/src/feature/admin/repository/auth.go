package repository

import (
	
	"telkom-property-service/src/model"
	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

type AuthRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAuthRepository(db *gorm.DB, logger *logrus.Logger) *AuthRepository {
	return &AuthRepository{db, logger}
}

func (r *AuthRepository) CreateAdmin(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AuthRepository) GetAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}