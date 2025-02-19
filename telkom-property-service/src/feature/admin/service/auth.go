package service

import (

	"errors"
	"time"

	"telkom-property-service/src/model"
	"telkom-property-service/src/feature/admin/repository"
	
	"golang.org/x/crypto/bcrypt"
	
	"github.com/sirupsen/logrus"
	"github.com/golang-jwt/jwt/v4"


)

type AuthService struct {
	repo   *repository.AuthRepository
	logger *logrus.Logger
}

func NewAuthService(repo *repository.AuthRepository, logger *logrus.Logger) *AuthService {
	return &AuthService{repo, logger}
}

func (s *AuthService) RegisterAdmin(username, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	
	admin := model.Admin{
		Username: username,
		Password: string(hashedPassword), 
	}

	return s.repo.CreateAdmin(&admin)
}

func (s *AuthService) Authenticate(username, password string) (*model.Admin, error) {
	admin, err := s.repo.GetAdminByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, err
	}

	return admin, nil
}

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	AdminID uint   `json:"admin_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(adminID uint, username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute)

	claims := &Claims{
		AdminID:  adminID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// func (s *AuthService) UpdateToken(AdminID uint, token string) error {
// 	admin := model.Admin{}
// 	if err := s.db.First(&admin, adminID).Error; err != nil {
// 		return err
// 	}

// 	admin.Token = token;
// 	return s.db.Save(&admin).Error
// }




