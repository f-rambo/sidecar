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

type SystemUsecase struct {
	log *log.Helper
}

func NewSystemUseCase(logger log.Logger) *SystemUsecase {
	s := &SystemUsecase{
		log: log.NewHelper(logger),
	}
	return s
}
