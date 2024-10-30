package interfaces

import (
	"context"

	v1alpha1 "github.com/f-rambo/ship/api/cloud/v1alpha1"
	"github.com/f-rambo/ship/api/common"
	"github.com/f-rambo/ship/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
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

// node init
func (c *CloudInterface) NodeInit(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	cloud := c.interfaceToBiz(req)
	if cloud.NodeID == 0 || cloud.NodeName == "" {
		return nil, errors.New("node id or node name is empty")
	}
	err := c.uc.NodeInit(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// InstallKubeadmKubeletCriO
func (c *CloudInterface) InstallKubeadmKubeletCriO(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	cloud := c.interfaceToBiz(req)
	if cloud.ClusterVersion == "" {
		return nil, errors.New("cluster version is empty")
	}
	err := c.uc.InstallKubeadmKubeletCriO(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// KubeadmJoin
func (c *CloudInterface) KubeadmJoin(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	if req.ControlPlaneEndpoint == "" {
		return nil, errors.New("control plane endpoint is empty")
	}
	if req.Token == "" {
		return nil, errors.New("token is empty")
	}
	if req.DiscoveryTokenCaCertHash == "" {
		return nil, errors.New("discovery token ca cert hash is empty")
	}
	cloud := c.interfaceToBiz(req)
	err := c.uc.KubeadmJoin(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// KubeadmInit
func (c *CloudInterface) KubeadmInit(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	cloud := c.interfaceToBiz(req)
	err := c.uc.KubeadmInit(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// KubeadmReset
func (c *CloudInterface) KubeadmReset(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	cloud := c.interfaceToBiz(req)
	err := c.uc.KubeadmReset(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// KubeadmUpgrade
func (c *CloudInterface) KubeadmUpgrade(ctx context.Context, req *v1alpha1.Cloud) (*common.Msg, error) {
	if req == nil {
		return nil, errors.New("cloud is empty")
	}
	cloud := c.interfaceToBiz(req)
	err := c.uc.KubeadmUpgrade(ctx, cloud)
	if err != nil {
		return nil, err
	}
	return common.Response(), nil
}

// interface cloud to biz cloud
func (c *CloudInterface) interfaceToBiz(cloud *v1alpha1.Cloud) *biz.Cloud {
	return &biz.Cloud{
		ID:                       cloud.Id,
		NodeID:                   cloud.NodeId,
		NodeName:                 cloud.NodeName,
		ClusterVersion:           cloud.ClusterVersion,
		Token:                    cloud.Token,
		DiscoveryTokenCACertHash: cloud.DiscoveryTokenCaCertHash,
		ControlPlaneEndpoint:     cloud.ControlPlaneEndpoint,
		JoinConfig:               cloud.JoinConfig,
	}
}
