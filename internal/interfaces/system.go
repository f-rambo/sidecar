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

func (c *SystemInterface) GetSystem(ctx context.Context, _ *emptypb.Empty) (*v1alpha1.System, error) {
	system, err := c.systemUc.GetSystem(ctx)
	if err != nil {
		return nil, err
	}
	return &v1alpha1.System{
		Id:         system.ID,
		Os:         system.OS,
		Arch:       system.ARCH,
		Cpu:        system.CPU,
		Memory:     system.Memory,
		Gpu:        system.GPU,
		GpuSpec:    system.GpuSpec,
		DataDisk:   system.DataDisk,
		Kernel:     system.Kernel,
		Container:  system.Container,
		Kubelet:    system.Kubelet,
		KubeProxy:  system.KubeProxy,
		InternalIp: system.InternalIP,
	}, nil
}
