package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorWorks(t *testing.T) {
	sessionId := generateSessionId()
	assert.NotEmpty(t, sessionId)
}
