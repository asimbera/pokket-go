package main

import (
	"bytes"
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

func TestAuth(t *testing.T) {
	r := setupRouter()

	t.Run("InvalidSignupForm", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "asimbera@outlook.in"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "Invalid request body", w.Body.String())
	})

	t.Run("SignupSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "asimbera@outlook.in", "name": "Asim Bera", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	t.Run("EmailAlreadyRegistered", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "asimbera@outlook.in", "name": "Asim Bera", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "Failed to create user", w.Body.String())
	})

	t.Run("InvalidLoginForm", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "asimbera@outlook.in"}`)
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
		reqBody := []byte(`{"email": "asimbera@outlook.in", "password": "1234567890"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Equal(t, "Invalid password", w.Body.String())
	})

	t.Run("LoginSuccess", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := []byte(`{"email": "asimbera@outlook.in", "password": "12345678"}`)
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	// Teardown
	var user models.User
	models.Database.Where("email = ?", "asimbera@outlook.in").Find(&user)
	models.Database.Unscoped().Delete(&user)

}
