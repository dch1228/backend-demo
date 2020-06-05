package tracing

import (
	"io"

	"github.com/duchenhao/backend-demo/internal/conf"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var closer io.Closer

func Init() {
	cfg := jaegercfg.Configuration{
		ServiceName: conf.Core.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}

	var tracer opentracing.Tracer
	var err error
	tracer, closer, err = cfg.NewTracer()
	if err != nil {
		panic(err)
	}

	opentracing.SetGlobalTracer(tracer)
}

func Close() {
	closer.Close()
}
