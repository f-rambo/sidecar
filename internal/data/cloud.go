package data

import (
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type cloudRepo struct {
	data *Data
	log  *log.Helper
}

func NewCloudRepo(data *Data, logger log.Logger) biz.CloudRepo {
	return &cloudRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
