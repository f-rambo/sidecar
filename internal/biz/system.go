package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "system.c"
*/
import "C"

import (
	"github.com/go-kratos/kratos/v2/log"
)

type System struct {
}

type SystemRepo interface {
}

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
