package controllers

import (
	"github.com/asimbera/pokket/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const JwtSecret = "my_secret_string"

type LoginForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignupForm struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


func LoginController(c *gin.Context) {
	// Validate request body
	form := &LoginForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}
	// Check user in database
	var user models.User
	if err := models.Database.Where(&models.User{Email: form.Email}).First(&user).Error; err != nil {
		c.String(http.StatusForbidden, "User not found")
		return
	}

	// Match user's password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		c.String(http.StatusForbidden, "Invalid password")
		return
	}
	// Sign JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to sign token")
		return
	}

	// Return JWT token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})
	return
}

func SignupController(c *gin.Context) {
	// Validate request body
	form := &SignupForm{}
	if err := c.ShouldBindJSON(form); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}
	// Check if email already exists in database
	//var count int64
	//if models.Database.Where(&models.User{Email: form.Email}).Count(&count); count >= 1 {
	//	c.String(http.StatusForbidden, "Email already registered")
	//	return
	//}
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), 10)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while hashing password")
		return
	}
	// Store user in database
	user := &models.User{
		Name: form.Name,
		Email:           form.Email,
		Password:        string(hash),
		IsEmailVerified: false,
	}
	if err := models.Database.Create(&user).Error; err != nil {
		//c.String(http.StatusInternalServerError, "Failed to create user")
		c.String(http.StatusForbidden, "Failed to create user")
		return
	}

	// Sign JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to sign token")
		return
	}

	// Return JWT token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
	})
	return
}
