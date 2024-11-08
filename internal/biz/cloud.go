package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "cloud.c"
*/
import "C"

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Cloud struct {
}

type CloudUsecase struct {
	log *log.Helper
}

func NewCloudUseCase(logger log.Logger) *CloudUsecase {
	c := &CloudUsecase{
		log: log.NewHelper(logger),
	}
	return c
}
