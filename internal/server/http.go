package server

import (
	"time"

	clusterv1alpha1 "github.com/f-rambo/ship/api/cluster/v1alpha1"
	systemv1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/internal/interfaces"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, system *interfaces.SystemInterface, cluster *interfaces.ClusterInterface, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.HTTP.Network != "" {
		opts = append(opts, http.Network(c.HTTP.Network))
	}
	if c.HTTP.Addr != "" {
		opts = append(opts, http.Address(c.HTTP.Addr))
	}
	if c.HTTP.Timeout != 0 {
		opts = append(opts, http.Timeout(time.Duration(c.GRPC.Timeout)*time.Second))
	}
	srv := http.NewServer(opts...)
	clusterv1alpha1.RegisterClusterInterfaceHTTPServer(srv, cluster)
	systemv1alpha1.RegisterSystemInterfaceHTTPServer(srv, system)
	return srv
}
