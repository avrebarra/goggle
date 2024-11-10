package domainaccesslog

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func (AccessLog) Fake() AccessLog {
	return AccessLog{
		ID:        gofakeit.Int(),
		ToggleID:  "general/" + gofakeit.UUID(),
		CreatedAt: gofakeit.DateRange(time.Now().Add(-1*30*24*time.Hour), time.Now()),
	}
}
