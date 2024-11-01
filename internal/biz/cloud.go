package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "cloud.c"
*/
import "C"

import (
	"github.com/f-rambo/ship/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
)

type Cloud struct {
}

type CloudRepo interface{}

type CloudUsecase struct {
	cloudRepo CloudRepo
	log       *log.Helper
	c         *conf.Server
}

func NewCloudUseCase(conf *conf.Server, cloudRepo CloudRepo, logger log.Logger) *CloudUsecase {
	c := &CloudUsecase{
		cloudRepo: cloudRepo,
		log:       log.NewHelper(logger),
		c:         conf,
	}
	return c
}
