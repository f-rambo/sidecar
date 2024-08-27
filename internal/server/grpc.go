package server

import (
	"time"

	clusterv1alpha1 "github.com/f-rambo/ship/api/cluster/v1alpha1"
	systemv1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/internal/interfaces"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, system *interfaces.SystemInterface, cluster *interfaces.ClusterInterface, logger log.Logger) *grpc.Server {
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
	clusterv1alpha1.RegisterClusterInterfaceServer(srv, cluster)
	systemv1alpha1.RegisterSystemInterfaceServer(srv, system)
	return srv
}
