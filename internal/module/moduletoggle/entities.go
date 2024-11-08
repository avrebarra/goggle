package moduletoggle

import (
	"time"
)

type Toggle struct {
	ID        string    `validate:"required"`
	Status    bool      `validate:"-"`
	UpdatedAt time.Time `validate:"required"`
}

type ToggleWithDetail struct {
	ID               string
	Status           bool
	UpdatedAt        time.Time
	LastAccessedAt   time.Time
	AccessFreqWeekly int
}

type ToggleStat struct {
	ID     string `validate:"required"`
	Status bool   `validate:"-"`
}
