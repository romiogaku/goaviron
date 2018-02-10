//go:generate goagen bootstrap -d goaviron/design

package main

import (
	"compress/gzip"
	"flag"
	"goaviron/app"
	"goaviron/env"
	"goaviron/presentation/controller"

	"github.com/deadcheat/goacors"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	gm "github.com/goadesign/goa/middleware/gzip"
)

func main() {
	// Create service
	service := goa.New("goaviron")

	// cors
	service.Use(goacors.WithConfig(service, &env.CorsConf))

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())
	service.Use(middleware.Recover())
	service.Use(gm.Middleware(gzip.BestSpeed))

	// Mount "swagger" controller
	swg := controller.NewSwaggerController(service)
	controller.MountSwaggerController(service, swg)

	// Mount "viron" controller
	v := controller.NewVironController(service)
	app.MountVironController(service, v)

	// Mount "aikatsu" controller
	a := controller.NewAikatsuController(service)
	app.MountAikatsuController(service, a)

	// コマンド引数で起動ポート・起動ホストを上書きできるようにする
	port := flag.Int("p", env.Server.PortNum, "port number. default set on config")
	docPort := flag.Int("dp", env.Server.DocPort, "port number for doc. default set on config")
	host := flag.String("h", env.Server.HostName, "name of server host. default set on config")
	docHost := flag.String("dh", env.Server.DocHostName, "name of server host. default set on config")
	flag.Parse()

	env.Server.HostName = *host
	env.Server.PortNum = *port

	env.Server.DocHostName = *docHost
	env.Server.DocPort = *docPort

	// Start service
	if err := service.ListenAndServe(env.Server.APIHostString()); err != nil {
		service.LogError("startup", "err", err)
	}

}
