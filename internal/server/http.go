package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duchenhao/backend-demo/internal/api"
	"github.com/duchenhao/backend-demo/internal/conf"
	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/middleware"
)

type HttpServer struct {
	log     *zap.Logger
	httpSrv *http.Server
	gin     *gin.Engine
}

func NewHttpServer() *HttpServer {
	hs := &HttpServer{
		log: log.Named("http.server"),
	}

	hs.newGin()
	hs.addMiddleware()
	api.RegisterRoutes(hs.gin)

	return hs
}

func (hs *HttpServer) Run() {
	hs.httpSrv = &http.Server{
		Addr:    conf.Core.Addr,
		Handler: hs.gin,
	}

	hs.log.Info("http server listen", zap.String("address", hs.httpSrv.Addr))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := hs.httpSrv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				hs.log.Info("server was shutdown gracefully")
				return
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	if err := hs.httpSrv.Shutdown(context.Background()); err != nil {
		hs.log.Error("Failed to shutdown server", zap.Error(err))
	}
	wg.Wait()
}

func (hs *HttpServer) newGin() {
	gin.SetMode(conf.GetHttpEnv())

	g := gin.New()

	hs.gin = g
}

func (hs *HttpServer) addMiddleware() {
	g := hs.gin

	g.Use(middleware.Logger())
	g.Use(middleware.Recovery())
	g.Use(middleware.RequestMetrics(g))
	g.Use(middleware.ContextHandler())
}
