package data

import (
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type clusterRepo struct {
	data *Data
	log  *log.Helper
}

func NewClusterRepo(data *Data, logger log.Logger) biz.ClusterRepo {
	return &clusterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
