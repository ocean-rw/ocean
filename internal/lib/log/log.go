package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
}

type Config struct {
	Level       zap.AtomicLevel `yaml:"level"`
	Development bool            `yaml:"development"`
	StdoutPaths []string        `yaml:"stdout_paths"`
	StderrPaths []string        `yaml:"stderr_paths"`
}

func New(cfg *Config) *zap.SugaredLogger {
	logCfg := zap.NewProductionConfig()
	logCfg.Sampling = nil
	if cfg != nil {
		logCfg.Development = cfg.Development
		if cfg.Level != (zap.AtomicLevel{}) {
			logCfg.Level = cfg.Level
		}
		if len(cfg.StdoutPaths) > 0 {
			logCfg.OutputPaths = cfg.StdoutPaths
		}
		if len(cfg.StderrPaths) > 0 {
			logCfg.ErrorOutputPaths = cfg.StderrPaths
		}
	}
	logger, _ := logCfg.Build()
	return logger.Sugar()
}

func Middleware(l *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			now := time.Now()
			defer func() {
				l.Infof("%s %s %s %s %d %d %d",
					middleware.GetReqID(r.Context()),
					r.Method, r.URL.Path, r.Proto,
					ww.Status(), ww.BytesWritten(), time.Since(now).Nanoseconds(),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
