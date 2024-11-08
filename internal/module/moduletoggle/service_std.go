package moduletoggle

import "github.com/avrebarra/goggle/utils/validator"

type ServiceConfig struct {
}

type ServiceStd struct {
	Config ServiceConfig
}

func NewService(cfg ServiceConfig) (out *ServiceStd, err error) {
	if err = validator.Validate(&cfg); err != nil {
		return nil, err
	}
	out = &ServiceStd{Config: cfg}
	return
}
