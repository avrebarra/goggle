package clientauth

import (
	"context"
	"strings"
)

var _ Client = (*ClientDummy)(nil)

type ClientDummy struct {
	// Mock responses for testing
	Users map[string]*UserInfo
}

func NewClientDummy() *ClientDummy {
	return &ClientDummy{
		Users: map[string]*UserInfo{
			"valid-token-123": {
				UserID:   "user-001",
				Username: "john_doe",
				Email:    "john@example.com",
			},
			"admin-token-456": {
				UserID:   "user-002",
				Username: "admin",
				Email:    "admin@example.com",
			},
		},
	}
}

func (c *ClientDummy) ValidateSession(ctx context.Context, token string) (*UserInfo, error) {
	// Simulate invalid token
	if token == "" || strings.Contains(token, "invalid") {
		return nil, ErrInvalidToken
	}

	// Simulate unauthorized token
	if strings.Contains(token, "unauthorized") {
		return nil, ErrUnauthorized
	}

	// Simulate service down
	if strings.Contains(token, "service-down") {
		return nil, ErrServiceDown
	}

	// Return mock user if token exists
	if user, exists := c.Users[token]; exists {
		return user, nil
	}

	// Default invalid token
	return nil, ErrInvalidToken
}
