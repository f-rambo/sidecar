package server

import (
	"time"

	cloudv1alpha1 "github.com/f-rambo/ship/api/cloud/v1alpha1"
	systemv1alpha1 "github.com/f-rambo/ship/api/system/v1alpha1"
	"github.com/f-rambo/ship/internal/conf"
	"github.com/f-rambo/ship/internal/interfaces"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, system *interfaces.SystemInterface, cloud *interfaces.CloudInterface, logger log.Logger) *http.Server {
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
	cloudv1alpha1.RegisterCloudInterfaceHTTPServer(srv, cloud)
	systemv1alpha1.RegisterSystemInterfaceHTTPServer(srv, system)
	return srv
}
