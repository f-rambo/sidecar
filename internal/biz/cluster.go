package biz

import "github.com/go-kratos/kratos/v2/log"

type ClusterRepo interface{}

type ClusterUsecase struct {
	clusterRepo ClusterRepo
	log         *log.Helper
}

func NewClusterUseCase(clusterRepo ClusterRepo, logger log.Logger) *ClusterUsecase {
	c := &ClusterUsecase{
		clusterRepo: clusterRepo,
		log:         log.NewHelper(logger),
	}
	return c
}
