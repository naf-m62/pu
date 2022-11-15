package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"pu/config"
)

type (
	Logger = *zap.Logger

	loggerConfig struct {
		Level       string `mapstructure:"level"`
		Development bool   `mapstructure:"development"`
		Caller      bool   `mapstructure:"caller"`
		Stacktrace  string `mapstructure:"stacktrace"`
	}
)

func New(config config.Config) (l Logger, err error) {
	lCfg := &loggerConfig{}
	if err = config.UnmarshalKey("logger", &lCfg); err != nil {
		return nil, err
	}

	var level zap.AtomicLevel
	if err = level.UnmarshalText([]byte(lCfg.Level)); err != nil {
		return nil, err
	}

	eConf := zap.NewDevelopmentEncoderConfig()
	eConf.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var options = make([]zap.Option, 0, 3)
	if lCfg.Caller {
		options = append(options, zap.AddCaller())
	}

	if lCfg.Development {
		options = append(options, zap.Development())
	}

	var sLevel zap.AtomicLevel
	if len(lCfg.Stacktrace) > 0 {
		if err = sLevel.UnmarshalText([]byte(lCfg.Stacktrace)); err != nil {
			return nil, err
		}

		options = append(options, zap.AddStacktrace(sLevel))
	}

	l = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(eConf),
		zapcore.AddSync(os.Stdout),
		level,
	), options...)

	return l, nil
}
