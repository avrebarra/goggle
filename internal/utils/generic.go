package utils

import (
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

func ApplyDefaults[T any](targ, defs *T) {
	var tmp T
	_ = copier.Copy(&tmp, defs)
	_ = copier.CopyWithOption(&tmp, targ, copier.Option{IgnoreEmpty: true})
	_ = copier.Copy(targ, &tmp)
}

func MorphFrom[T any](targ *T, vals any, supplements *T) (err error) {
	if err = copier.Copy(targ, vals); err != nil {
		err = errors.Wrap(err, "failed to compose")
		return
	}
	if supplements != nil {
		if err = copier.CopyWithOption(targ, supplements, copier.Option{IgnoreEmpty: true}); err != nil {
			err = errors.Wrap(err, "failed to supplement")
			return
		}
	}
	return
}

func UnmarshalToMap(bytes []byte) (m map[string]interface{}) {
	m = map[string]interface{}{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return
	}
	return
}

func RemarshalToMap(i interface{}) (m map[string]interface{}) {
	m = map[string]interface{}{}
	bytes, err := json.Marshal(i)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return
	}
	return
}
