package interfaces

import (
	"context"

	v1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SystemInterface struct {
	v1alpha1.UnimplementedSystemInterfaceServer
	systemUc *biz.SystemUsecase
	log      *log.Helper
}

func NewSystemInterface(systemUc *biz.SystemUsecase, logger log.Logger) *SystemInterface {
	return &SystemInterface{
		systemUc: systemUc,
		log:      log.NewHelper(logger),
	}
}

func (c *SystemInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*v1alpha1.Msg, error) {
	return &v1alpha1.Msg{Message: "pong"}, nil
}
