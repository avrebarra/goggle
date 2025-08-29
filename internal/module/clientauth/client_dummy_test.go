package clientauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientDummy_ValidateSession(t *testing.T) {
	client := NewClientDummy()
	ctx := context.Background()

	t.Run("valid token", func(t *testing.T) {
		user, err := client.ValidateSession(ctx, "valid-token-123")
		require.NoError(t, err)
		assert.Equal(t, "user-001", user.UserID)
		assert.Equal(t, "john_doe", user.Username)
		assert.Equal(t, "john@example.com", user.Email)
	})

	t.Run("admin token", func(t *testing.T) {
		user, err := client.ValidateSession(ctx, "admin-token-456")
		require.NoError(t, err)
		assert.Equal(t, "user-002", user.UserID)
		assert.Equal(t, "admin", user.Username)
		assert.Equal(t, "admin@example.com", user.Email)
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := client.ValidateSession(ctx, "")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := client.ValidateSession(ctx, "invalid-token")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("unauthorized token", func(t *testing.T) {
		_, err := client.ValidateSession(ctx, "unauthorized-token")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("service down", func(t *testing.T) {
		_, err := client.ValidateSession(ctx, "service-down-token")
		assert.ErrorIs(t, err, ErrServiceDown)
	})

	t.Run("unknown token", func(t *testing.T) {
		_, err := client.ValidateSession(ctx, "unknown-token")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})
}
