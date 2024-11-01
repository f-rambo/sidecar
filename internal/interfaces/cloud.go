package interfaces

import (
	"context"

	v1alpha1 "github.com/f-rambo/ship/api/cloud/v1alpha1"
	"github.com/f-rambo/ship/api/common"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CloudInterface struct {
	v1alpha1.UnimplementedCloudInterfaceServer
	uc  *biz.CloudUsecase
	log *log.Helper
}

func NewCloudInterface(uc *biz.CloudUsecase, logger log.Logger) *CloudInterface {
	return &CloudInterface{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (c *CloudInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*common.Msg, error) {
	return common.Response(), nil
}
