package biz

import "github.com/go-kratos/kratos/v2/log"

type SystemRepo interface{}

type SystemUsecase struct {
	SystemRepo SystemRepo
	log        *log.Helper
}

func NewSystemUseCase(systemRepo SystemRepo, logger log.Logger) *SystemUsecase {
	s := &SystemUsecase{
		SystemRepo: systemRepo,
		log:        log.NewHelper(logger),
	}
	return s
}
