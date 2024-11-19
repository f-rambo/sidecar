package biz

/*
#cgo CFLAGS: -I../unix
#cgo LDFLAGS:
#include "cluster.c"
*/
import "C"

import (
	"github.com/go-kratos/kratos/v2/log"
)

type Cluster struct {
}

type ClusterUsecase struct {
	log *log.Helper
}

func NewClusterUseCase(logger log.Logger) *ClusterUsecase {
	c := &ClusterUsecase{
		log: log.NewHelper(logger),
	}
	return c
}
