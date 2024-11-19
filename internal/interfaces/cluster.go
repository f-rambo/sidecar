package interfaces

import (
	"context"

	clusterApi "github.com/f-rambo/cloud-copilot/sidecar/api/cluster"
	"github.com/f-rambo/cloud-copilot/sidecar/api/common"
	"github.com/f-rambo/cloud-copilot/sidecar/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterInterface struct {
	clusterApi.UnimplementedClusterInterfaceServer
	uc  *biz.ClusterUsecase
	log *log.Helper
}

func NewClusterInterface(uc *biz.ClusterUsecase, logger log.Logger) *ClusterInterface {
	return &ClusterInterface{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (c *ClusterInterface) Ping(ctx context.Context, _ *emptypb.Empty) (*common.Msg, error) {
	return common.Response(), nil
}

func (c *ClusterInterface) Info(ctx context.Context, _ *emptypb.Empty) (*clusterApi.Cluster, error) {
	return &clusterApi.Cluster{}, nil
}
