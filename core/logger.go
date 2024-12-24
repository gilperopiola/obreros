package core

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/obreros/core/errs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

// We use zap. It's fast and easy.
// Set it up and then just use it with zap.L() or zap.S().

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Logger -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

const LogsTimeLayout = "02/01/06 15:04:05"

// Replaces the global Logger in the zap pkg with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(cfg *LoggerCfg) *zap.Logger {
	zapOpts := newZapBuildOpts(cfg.LevelStackT)

	zapLogger, err := newZapConfig(cfg).Build(zapOpts...)
	if err != nil {
		log.Fatalf(errs.FailedToCreateLogger, err) // Don't use zap for this.
	}

	zap.ReplaceGlobals(zapLogger)

	return zapLogger
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Helpers -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type logOption func(*zap.Logger)

// Prepares a new child Logger with fields defined by the logOpts.
func willLog(opts ...logOption) *zap.Logger {
	l := zap.L()
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Prefix used when Infof or Infoln are called.
var WorkerLogsPrefix = AppEmoji + " " + AppAlias + " | "

func WorkerLog(s string) {
	zap.L().Info(WorkerLogsPrefix + s)
}

func WorkerLogf(s string, args ...any) {
	zap.S().Infof(WorkerLogsPrefix+s, args...)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Used to log unexpected errors, like panic recoveries or some connection errors.
func LogUnexpectedError(err error) {
	willLog(withError(err), withStacktrace()).Error("Unexpected Error 🛑")
}

// Used to log unexpected errors that also should trigger a panic.
func LogFatal(err error) {
	willLog(withError(err), withStacktrace()).Fatal("Unexpected Fatal Error 🛑")
}

// Helps keeping code clean and readable, lets you omit the error check on the caller when you just need to log it.
func WarnIfErr(err error, optionalFmt ...string) {
	if err != nil {
		format := "untyped warning: %v"
		if len(optionalFmt) > 0 {
			format = optionalFmt[0]
		}
		zap.S().Warnf(format, err)
	}
}

// Helps keeping code clean and readable, lets you omit the error check on the caller when you just need to log it.
// Use this for errors that are expected and handled.
func LogIfErr(err error, optionalFmt ...string) {
	if err != nil {
		format := "untyped error: %v"
		if len(optionalFmt) > 0 {
			format = optionalFmt[0]
		}
		zap.S().Errorf(format, err)
	}
}

// Helps keeping code clean and readable, lets you omit the error check on the caller.
func LogFatalIfErr(err error, optionalFormat ...string) {
	if err == nil {
		return
	}

	format := "untyped fatal: %v"
	if len(optionalFormat) > 0 {
		format = optionalFormat[0]
	}

	LogFatal(fmt.Errorf(format, err))
}

func LogImportant(msg string) {
	willLog(withMsg(msg)).Info("Important! ⭐")
}

func LogResult(resultOfWhat string, err error) {
	if err == nil {
		willLog().Info("✅ " + resultOfWhat + " succeeded!")
	} else {
		willLog(withError(err)).Error("❌ " + resultOfWhat + " failed!")
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - With Fields -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Logs a simple message.
var withMsg = func(msg string) logOption {
	return func(l *zap.Logger) {
		*l = *l.With(zap.String("msg", msg))
	}
}

// Logs any kind of info.
var withData = func(data ...any) logOption {
	return func(l *zap.Logger) {
		if len(data) == 0 {
			return
		}
		*l = *l.With(zap.Any("data", data))
	}
}

// Logs a duration.
var withDuration = func(duration time.Duration) logOption {
	return func(l *zap.Logger) {
		*l = *l.With(zap.Duration("duration", duration))
	}
}

// Log error if not nil.
var withError = func(err error) logOption {
	return func(l *zap.Logger) {
		if err == nil {
			return
		}
		*l = *l.With(zap.Error(err))
	}
}

// Used to log where in the code a message comes from.
var withStacktrace = func() logOption {
	return func(l *zap.Logger) {
		*l = *l.With(zap.Stack("trace"))
	}
}

// -> In HTTP, we join Method and Path -> 'GET /users'.
var withHTTPRoute = func(req *http.Request) logOption {
	return func(l *zap.Logger) {
		*l = *l.With(zap.String("route", req.Method+" "+req.URL.Path))
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Config -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the default options for creating the zap Logger.
func newZapBuildOpts(levelStackTrace int) []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zapcore.Level(levelStackTrace)),
		zap.WithClock(zapcore.DefaultClock),
	}
}

// Returns a new zap.Config with the default options + *LoggerCfg settings.
func newZapConfig(cfg *LoggerCfg) zap.Config {
	zapCfg := newZapBaseConfig()

	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zapCfg.DisableCaller = !cfg.LogCaller
	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(LogsTimeLayout))
	}
	zapCfg.EncoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(d.Truncate(time.Millisecond).String())
	}

	return zapCfg
}

// Returns the default zap.Config for the current environment.
func newZapBaseConfig() zap.Config {
	return zap.NewProductionConfig()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Only log messages with a level equal or higher than the one we set in the config.
var LogLevels = map[string]int{
	"debug":  int(zap.DebugLevel),
	"info":   int(zap.InfoLevel),
	"warn":   int(zap.WarnLevel),
	"error":  int(zap.ErrorLevel),
	"dpanic": int(zap.DPanicLevel),
	"panic":  int(zap.PanicLevel),
	"fatal":  int(zap.FatalLevel),
}

// The selected DB Log Level will be used to log all SQL queries.
// 'silent' disables all logs, 'error' will only log errors, 'warn' logs errors and warnings, and 'info' logs everything.
var DBLogLevels = map[string]int{
	"silent": int(gormLogger.Silent),
	"error":  int(gormLogger.Error),
	"warn":   int(gormLogger.Warn),
	"info":   int(gormLogger.Info),
}
