package main

import (
	"bytes"
	"encoding/json"
	"github.com/asimbera/pokket/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRoute(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/status", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Ok", w.Body.String())
}

func TestAuthRoutes(t *testing.T) {
	r := setupRouter()

	t.Run("InvalidSignupForm", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "Invalid request body", w.Body.String())
	})

	t.Run("SignupSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com", "name": "John Doe", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("EmailAlreadyRegistered", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com", "name": "John Doe", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "Failed to create user", w.Body.String())
	})

	t.Run("InvalidLoginForm", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "Invalid request body", w.Body.String())
	})

	t.Run("InvalidLoginEmail", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "non_existant@example.com", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "User not found", w.Body.String())
	})

	t.Run("InvalidLoginPassword", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com", "password": "1234567890"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "Invalid password", w.Body.String())
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "johndoe@example.com", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	// Teardown
	var user models.User
	models.Database.Where("email = ?", "johndoe@example.com").Find(&user)
	models.Database.Unscoped().Delete(&user)

}

func TestSecureRoutes(t *testing.T) {
	// Setup
	r := setupRouter()

	w := httptest.NewRecorder()
	reqBody := []byte(`{"email": "johndoe@example.com", "name": "John Doe", "password": "12345678"}`)
	req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct{Token string `json:"token"`}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	t.Run("TokenIsNotEmpty", func(t *testing.T) {
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/secure/me", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
	})

	t.Run("Authorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/secure/me", nil)
		req.Header.Set("X-Token", resp.Token)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	// Teardown
	var user models.User
	models.Database.Where("email = ?", "johndoe@example.com").Find(&user)
	models.Database.Unscoped().Delete(&user)
}