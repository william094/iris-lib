package iris_lib

import (
	"github.com/iris-contrib/middleware/prometheus"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Prometheus(app iris.Application, serverName string) {
	m := prometheus.New(serverName, prometheus.DefaultBuckets...)
	app.Use(m.ServeHTTP)
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))
}
