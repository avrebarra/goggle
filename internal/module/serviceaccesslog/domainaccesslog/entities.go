package domainaccesslog

import (
	"time"
)

type AccessLog struct {
	ID        int       `validate:"required"`
	ToggleID  string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
}
