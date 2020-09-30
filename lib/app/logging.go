package app

import (
	"fmt"
	"io"
	"os"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/KyberNetwork/cclog/lib/client"
	"github.com/KyberNetwork/reserve-data/cmd/mode"
)

const (
	modeFlag           = "mode"
	sentryDSNFlag      = "sentry-dsn"
	sentryLevelFlag    = "sentry-lv"
	defaultSentryLevel = "error"
	errorLevel         = "error"
	warnLevel          = "warn"
	fatalLevel         = "fatal"
	infoLevel          = "info"

	logOutputFlag = "log-output"
	// defaultLogOutput = "./log.log"
	ccLogAddr = "cclog-addr"
	cclogName = "cclog-name"
)

// NewSentryFlags returns flags to init sentry client
func NewSentryFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   sentryDSNFlag,
			EnvVar: "SENTRY_DSN",
			Usage:  "dsn for sentry client",
		}, cli.StringFlag{
			Name:   sentryLevelFlag,
			EnvVar: "SENTRY_LEVEL",
			Usage:  "log level report message to sentry (info, error, warn, fatal)",
			Value:  defaultSentryLevel,
		}, cli.StringFlag{
			Name:   logOutputFlag,
			EnvVar: "LOG_OUTPUT",
			Usage:  "log path to log file",
		},
		cli.StringFlag{
			Name:   ccLogAddr,
			Usage:  "cclog-address",
			Value:  "",
			EnvVar: "CCLOG_ADDR",
		}, cli.StringFlag{
			Name:   cclogName,
			Usage:  "cclog-name",
			Value:  "cex-account-data",
			EnvVar: "CCLOG_NAME",
		},
	}
}

type syncer interface {
	Sync() error
}

// NewFlusher creates a new syncer from given syncer that log a error message if failed to sync.
func NewFlusher(s syncer) func() {
	return func() {
		// ignore the error as the sync function will always fail in Linux
		// https://github.com/uber-go/zap/issues/370
		_ = s.Sync()
	}
}

func configRotateLog(outputFile string) io.Writer {
	var mw io.Writer
	if outputFile != "" {
		return io.MultiWriter(os.Stdout)
	}
	logger := &lumberjack.Logger{
		Filename:   outputFile,
		MaxBackups: 0,
		MaxAge:     0,
	}
	mw = io.MultiWriter(os.Stdout, logger)

	c := cron.New()
	err := c.AddFunc("@daily", func() {
		if lErr := logger.Rotate(); lErr != nil {
			panic(fmt.Sprintf("failed to rotate log, error: %s", lErr))
		}
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create log rotate daily: %s", err))
	}
	c.Start()
	return mw
}

// NewLogger creates a new logger instance.
// The type of logger instance will be different with different application running modes.
func NewLogger(c *cli.Context) (*zap.Logger, error) {
	var (
		l       *zap.Logger
		encoder zapcore.Encoder
		writers = []io.Writer{os.Stdout}
	)
	logOutPutFile := c.String(logOutputFlag)
	logAddr := c.GlobalString(ccLogAddr)
	logName := c.GlobalString(cclogName)
	if logAddr != "" && logName != "" {
		ccw := client.NewAsyncLogClient(logName, logAddr, func(err error) {
			fmt.Println("send log error", err)
		})
		writers = append(writers, &UnescapeWriter{w: ccw})
		w := configRotateLog(logOutPutFile)
		writers = append(writers, w)
	}
	w := io.MultiWriter(writers...)

	modeConfig := c.GlobalString(modeFlag)
	switch modeConfig {
	case mode.Production.String():
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	case mode.Development.String():
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	default:
		return nil, fmt.Errorf("invalid running mode: %q", modeConfig)
	}
	l = zap.New(zapcore.NewCore(encoder, zapcore.AddSync(w), zap.DebugLevel))

	return l, nil
}

// SentryDSNFromFlag return sentry dsn string from flag input
func SentryDSNFromFlag(c *cli.Context) string {
	return c.String(sentryDSNFlag)
}

// NewSugaredLogger creates a new sugared logger and a flush function. The flush function should be
// called by consumer before quitting application.
// This function should be use most of the time unless
// the application requires extensive performance, in this case use NewLogger.
func NewSugaredLogger(c *cli.Context) (*zap.SugaredLogger, func(), error) {
	logger, err := NewLogger(c)
	if err != nil {
		return nil, nil, err
	}
	// init sentry if flag dsn exists
	if len(c.String(sentryDSNFlag)) != 0 {
		sentryClient, err := sentry.NewClient(sentry.ClientOptions{
			Dsn: c.String(sentryDSNFlag),
		})
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to init sentry client")
		}

		cfg := zapsentry.Configuration{
			DisableStacktrace: false,
		}
		switch c.String(sentryLevelFlag) {
		case infoLevel:
			cfg.Level = zapcore.InfoLevel
		case warnLevel:
			cfg.Level = zapcore.WarnLevel
		case errorLevel:
			cfg.Level = zapcore.ErrorLevel
		case fatalLevel:
			cfg.Level = zapcore.FatalLevel
		default:
			return nil, nil, errors.Errorf("invalid log level %v", c.String(sentryLevelFlag))
		}

		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentryClient))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to init zap sentry")
		}
		// attach to logger core
		logger = zapsentry.AttachCoreToLogger(core, logger)
	}
	sugar := logger.Sugar()
	return sugar, NewFlusher(logger), nil
}
