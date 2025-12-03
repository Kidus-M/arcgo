package infrastructure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"task_manager1/Infrastructure/security"
)

func TestPasswordHashing(t *testing.T) {
	pw := security.NewPasswordService()

	hash, err := pw.HashPassword("hello")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "hello", hash)

	assert.True(t, pw.ComparePassword(hash, "hello"))
	assert.False(t, pw.ComparePassword(hash, "wrong"))
}
