package main

import (
	"errors"
	"flag"
	"os"

	"github.com/cristalhq/aconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	var (
		lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	)
	if os.Getenv("DEBUG") == "1" {
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	l, err := zap.Config{
		Level:       lvl,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		// stdout for result, logs must be present in stderr, reference:
		// https://docs.ansible.com/ansible/latest/collections/ansible/builtin/script_inventory.html#parameter-always_show_stderr
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}.Build()
	if err != nil {
		panic(err)
	}
	logger = l

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		// todo : boolean flags must not require arguments
		SkipFlags: true,
		EnvPrefix: "NSPAWN_INVENTORY",
	})
	loader.WalkFields(func(f aconfig.Field) bool {
		logger.Debug("field discovered", zap.String("name", f.Name()), zap.String("env", f.Tag("env")))
		return true
	})
	if err := loader.Load(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		logger.Fatal("failed to load configuration", zap.Error(err))
	}
}
