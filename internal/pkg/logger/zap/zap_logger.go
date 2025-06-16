package zap

import (
	"time"

	"github.com/DavidReque/go-food-delivery/internal/pkg/config/environment"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger"
	"github.com/DavidReque/go-food-delivery/internal/pkg/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapLogger(options *config.LogOptions, env environment.Environment) logger.Logger {
	// Configura el logger según el ambiente (development/production)
	var config zap.Config

	if env.IsProduction() {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if options != nil && !options.CallerEnabled {
		config.DisableCaller = true
	}

	zapLog, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		logger: zapLog,
		sugar:  zapLog.Sugar(),
	}
}

func (l *zapLogger) Configure(cfg func(interface{})) {
	if cfg != nil {
		cfg(l.logger)
	}
}

// Debug registra mensajes de nivel de depuración
func (l *zapLogger) Debug(args ...interface{})                   { l.sugar.Debug(args...) }
func (l *zapLogger) Debugf(template string, args ...interface{}) { l.sugar.Debugf(template, args...) }
func (l *zapLogger) Debugw(msg string, fields logger.Fields) {
	l.sugar.Debugw(msg, fieldsToArgs(fields)...)
}

// Info registra mensajes de nivel de información
func (l *zapLogger) Info(args ...interface{})                   { l.sugar.Info(args...) }
func (l *zapLogger) Infof(template string, args ...interface{}) { l.sugar.Infof(template, args...) }
func (l *zapLogger) Infow(msg string, fields logger.Fields) {
	l.sugar.Infow(msg, fieldsToArgs(fields)...)
}

// Warn registra mensajes de nivel de advertencia
func (l *zapLogger) Warn(args ...interface{})                   { l.sugar.Warn(args...) }
func (l *zapLogger) Warnf(template string, args ...interface{}) { l.sugar.Warnf(template, args...) }
func (l *zapLogger) WarnMsg(msg string, err error)              { l.sugar.Warnw(msg, "error", err) }

// Error registra mensajes de nivel de error
func (l *zapLogger) Error(args ...interface{}) { l.sugar.Error(args...) }
func (l *zapLogger) Errorw(msg string, fields logger.Fields) {
	l.sugar.Errorw(msg, fieldsToArgs(fields)...)
}
func (l *zapLogger) Errorf(template string, args ...interface{}) { l.sugar.Errorf(template, args...) }
func (l *zapLogger) Err(msg string, err error)                   { l.sugar.Errorw(msg, "error", err) }

// Fatal registra mensajes de nivel de error fatal
func (l *zapLogger) Fatal(args ...interface{})                   { l.sugar.Fatal(args...) }
func (l *zapLogger) Fatalf(template string, args ...interface{}) { l.sugar.Fatalf(template, args...) }

// Printf registra mensajes de nivel de información
func (l *zapLogger) Printf(template string, args ...interface{}) { l.sugar.Infof(template, args...) }

// WithName agrega un nombre al logger
func (l *zapLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugar = l.logger.Sugar()
}

// GrpcMiddlewareAccessLogger registra información de acceso para middleware gRPC
func (l *zapLogger) GrpcMiddlewareAccessLogger(method string, duration time.Duration, metadata map[string][]string, err error) {
	l.sugar.Infow("gRPC Access Log",
		"method", method,
		"duration", duration,
		"metadata", metadata,
		"error", err,
	)
}

// GrpcClientInterceptorLogger registra información para clientes gRPC
func (l *zapLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, duration time.Duration, metadata map[string][]string, err error) {
	l.sugar.Infow("gRPC Client Log",
		"method", method,
		"request", req,
		"reply", reply,
		"duration", duration,
		"metadata", metadata,
		"error", err,
	)
}

// fieldsToArgs convierte un mapa de campos en argumentos para la función zap.SugaredLogger
func fieldsToArgs(fields logger.Fields) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}
