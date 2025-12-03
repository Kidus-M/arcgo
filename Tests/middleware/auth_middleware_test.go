package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"task_manager1/Infrastructure/auth"
)

func TestAuthMiddlewareRejectsNoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// FIX: Removed unused variable 'jwtSvc'
	mw := auth.NewAuthMiddleware()
	r := gin.Default()

	// NOTE: Ensure 'AuthRequired()' is implemented in your 'auth.AuthMiddleware' struct
	r.GET("/protected", mw.Handle(), func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}