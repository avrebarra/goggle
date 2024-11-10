package domaintoggle

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func (ToggleWithDetail) Fake() ToggleWithDetail {
	e := ToggleWithDetail{
		ID:               strings.ToLower(fmt.Sprintf("general/%s_%s", gofakeit.AdjectiveDescriptive(), gofakeit.Animal())),
		Status:           gofakeit.Bool(),
		LastAccessedAt:   time.Time{},
		UpdatedAt:        time.Time{},
		AccessFreqWeekly: 0,
	}
	if gofakeit.Bool() {
		e.AccessFreqWeekly = gofakeit.Number(20, 100)
		e.LastAccessedAt = gofakeit.DateRange(time.Now().Add(-1*30*24*time.Hour), time.Now())
		e.UpdatedAt = gofakeit.DateRange(time.Now().Add(-1*30*24*time.Hour), e.LastAccessedAt)
	}
	return e
}
