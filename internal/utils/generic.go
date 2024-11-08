package utils

import (
	"fmt"

	"github.com/jinzhu/copier"
)

func ApplyDefaults[T any](targ, defs *T) {
	var tmp T
	_ = copier.Copy(&tmp, defs)
	_ = copier.CopyWithOption(&tmp, targ, copier.Option{IgnoreEmpty: true})
	_ = copier.Copy(targ, &tmp)
}

func Translate[T any](targ *T, vals any, supplements *T) (err error) {
	if err = copier.Copy(targ, vals); err != nil {
		err = fmt.Errorf("failed to compose: %w", err)
		return
	}
	if supplements != nil {
		if err = copier.CopyWithOption(targ, supplements, copier.Option{IgnoreEmpty: true}); err != nil {
			err = fmt.Errorf("failed to supplement: %w", err)
			return
		}
	}
	return
}
