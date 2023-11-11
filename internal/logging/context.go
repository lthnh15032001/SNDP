package logging

import (
	"context"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/lthnh15032001/ngrok-impl/internal/constants"

	log "github.com/sirupsen/logrus"
)

type loggerContextKey struct{}

var globalLogger *log.Entry

func init() {
	globalLogger = log.New().WithFields(nil)
	level, err := log.ParseLevel(GetEnv(constants.ENV_LOG_LEVEL, "INFO"))
	if err != nil {
		panic(err)
	}
	globalLogger.Logger.SetLevel(level)
	globalLogger.Logger.SetReportCaller(true)
	globalLogger.Logger.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})
}

// ContextWithLogger returns a context.Context that has been augmented with
// the provided log.Entry.
func ContextWithLogger(ctx context.Context, logger *log.Entry) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// LoggerFromContext extracts a *log.Entry from the provided context.Context and
// returns it. If no *log.Entry is found, a global, error-level *log.Entry is
// returned.
func LoggerFromContext(ctx context.Context) *log.Entry {
	if logger := ctx.Value(loggerContextKey{}); logger != nil {
		return ctx.Value(loggerContextKey{}).(*log.Entry) // nolint: forcetypeassert
	}
	return globalLogger
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
