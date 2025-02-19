package controller

import (
	"net/http"
	"telkom-property-service/src/model"
	"telkom-property-service/src/feature/admin/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{service}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var admin model.Admin 

	if err := c.ShouldBindJSON(&admin); err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	
	if err := ctrl.service.RegisterAdmin(admin.Username, admin.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"admin":   admin, 
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var credentials model.Admin 

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	admin, err := ctrl.service.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}


	token, err := service.GenerateJWT(admin.ID, admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}


	// if err := ctrl.service.UpdateToken(admin.ID, token); 
	// 	err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save token"})
	// 		return
	// 	}


	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"admin": gin.H{
			"username": admin.Username,
			"id":       admin.ID,
			"token": token,
		},
	})

}

