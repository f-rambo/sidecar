package data

import (
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type systemRepo struct {
	data *Data
	log  *log.Helper
}

func NewSystemRepo(data *Data, logger log.Logger) biz.SystemRepo {
	return &clusterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
