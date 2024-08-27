package interfaces

import (
	"context"

	v1alpha1 "github.com/f-rambo/ship/api/cluster/v1alpha1"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterInterface struct {
	v1alpha1.UnimplementedClusterInterfaceServer
	clusterUc *biz.ClusterUsecase
	log       *log.Helper
}

func NewClusterInterface(clusterUc *biz.ClusterUsecase, logger log.Logger) *ClusterInterface {
	return &ClusterInterface{
		clusterUc: clusterUc,
		log:       log.NewHelper(logger),
	}
}

func (c *ClusterInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*v1alpha1.Msg, error) {
	return &v1alpha1.Msg{Message: "pong"}, nil
}
