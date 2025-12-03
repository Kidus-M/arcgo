package infrastructure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"task_manager1/Infrastructure/auth"
)

func TestJWTGenerationValidation(t *testing.T) {
	svc := auth.NewJWTService("test-secret")

	token, err := svc.GenerateToken("kidus", "admin")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "kidus", claims.Username)
	assert.Equal(t, "admin", claims.Role)
}
