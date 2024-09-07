package server

import (
	"time"

	cloudv1alpha1 "github.com/f-rambo/ship/api/cloud/v1alpha1"
	systemv1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/internal/interfaces"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, system *interfaces.SystemInterface, cloud *interfaces.CloudInterface, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.GRPC.Network != "" {
		opts = append(opts, grpc.Network(c.GRPC.Network))
	}
	if c.GRPC.Addr != "" {
		opts = append(opts, grpc.Address(c.GRPC.Addr))
	}
	if c.GRPC.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(c.GRPC.Timeout)*time.Second))
	}
	srv := grpc.NewServer(opts...)
	cloudv1alpha1.RegisterCloudInterfaceServer(srv, cloud)
	systemv1alpha1.RegisterSystemInterfaceServer(srv, system)
	return srv
}
