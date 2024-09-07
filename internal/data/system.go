package data

import (
	"context"

	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type systemRepo struct {
	data *Data
	log  *log.Helper
}

func NewSystemRepo(data *Data, logger log.Logger) biz.SystemRepo {
	return &systemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s *systemRepo) GetSystem(ctx context.Context) (*biz.System, error) {
	system := &biz.System{ID: 1}
	err := s.data.db.First(system).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return system, nil
}

func (s *systemRepo) SaveSystem(ctx context.Context, system *biz.System) error {
	return s.data.db.Save(system).Error
}
