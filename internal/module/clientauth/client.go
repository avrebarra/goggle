package clientauth

import (
	"context"
	"fmt"
)

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrInvalidToken = fmt.Errorf("invalid token")
	ErrServiceDown  = fmt.Errorf("auth service unavailable")
)

// Client interface for auth service integration
type Client interface {
	// ValidateSession translates session token into user information
	ValidateSession(ctx context.Context, token string) (user *UserInfo, err error)
}

// UserInfo represents user information extracted from session token
type UserInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
