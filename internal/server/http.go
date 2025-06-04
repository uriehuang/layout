package server

import (
	"layout/internal/conf"
	"layout/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	pHttp "github.com/uriehuang/pkg/http"
	pkgMetrics "github.com/uriehuang/pkg/metrics"
	"github.com/uriehuang/pkg/middleware/ctx"
	"github.com/uriehuang/pkg/middleware/metrics"
	"github.com/uriehuang/pkg/middleware/sign"
	v1 "github.com/uriehuang/protocol/api/helloworld/v1"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metrics.HttpServer(
				"layout",
				c.GetEnv(),
				metrics.WithMillSeconds(pkgMetrics.HttpServerRequestMillSecondsHistogram),
				metrics.WithRequests(pkgMetrics.HttpServerRequestsCounter),
			),
			ctx.Ctx(),
			sign.CheckSign(),
		),
		http.ErrorEncoder(pHttp.ErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
