package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	ConsoleEncoder = "console"
	JSONEncoder    = "json"
)

type Options func(*ZapConfig)

type ZapConfig struct {
	encoding          string
	level             string
	development       bool
	timeKey           string
	levelKey          string
	nameKey           string
	callerKey         string
	messageKey        string
	stackTraceKey     string
	lineEnding        string
	initialFields     map[string]interface{}
	outputPaths       []string
	callSkip          int
	timeEncoderLayout string
}

func WithLogEncoding(encoding string) Options {
	return func(c *ZapConfig) {
		c.encoding = encoding
		if encoding == "" {
			c.encoding = "json"
		}
	}
}

func WithLogLevel(level string) Options {
	return func(c *ZapConfig) {
		c.level = level
		if level == "" {
			c.level = "debug"
		}
	}
}

func WithLogDevelopmentMode(mode bool) Options {
	return func(c *ZapConfig) {
		c.development = mode
	}
}

func WithLogTimeKey(timeKey string) Options {
	return func(c *ZapConfig) {
		c.timeKey = timeKey
		if timeKey != "" {
			c.timeKey = "time"
		}
	}
}

func WithLogLevelKey(levelKey string) Options {
	return func(c *ZapConfig) {
		c.levelKey = levelKey
		if c.levelKey == "" {
			c.levelKey = "level"
		}
	}
}

func WithLogNameKey(nameKey string) Options {
	return func(c *ZapConfig) {
		c.nameKey = nameKey
		if c.nameKey == "" {
			c.nameKey = "name"
		}
	}
}

func WithLogCallerKey(callerKey string) Options {
	return func(c *ZapConfig) {
		c.callerKey = callerKey
		if callerKey == "" {
			c.callerKey = "caller"
		}
	}
}

func WithMessageKey(messageKey string) Options {
	return func(c *ZapConfig) {
		c.messageKey = messageKey
		if messageKey == "" {
			c.messageKey = "msg"
		}
	}
}

func WithLogStackTraceKey(stackTraceKey string) Options {
	return func(c *ZapConfig) {
		c.stackTraceKey = stackTraceKey
		if stackTraceKey == "" {
			c.stackTraceKey = "stack"
		}
	}
}

func WithLogLineEnding(lineEnding string) Options {
	return func(c *ZapConfig) {
		c.lineEnding = lineEnding
		if lineEnding == "" {
			c.lineEnding = zapcore.DefaultLineEnding
		}
	}
}

func WithLogInitialFields(initialFields map[string]interface{}) Options {
	return func(c *ZapConfig) {
		c.initialFields = initialFields
	}
}

func WithLogOutputPaths(outputPath []string) Options {
	return func(c *ZapConfig) {
		c.outputPaths = outputPath
		if len(c.outputPaths) == 0 {
			c.outputPaths = []string{"stdout"}
		}
	}
}

func WithLogCallerSkip(callSkip int) Options {
	return func(c *ZapConfig) {
		c.callSkip = callSkip
		if callSkip < 1 {
			c.callSkip = 1
		}
	}
}

func WithTimeEncoderLayout(layout string) Options {
	return func(c *ZapConfig) {
		c.timeEncoderLayout = layout
		if layout == "" {
			c.timeEncoderLayout = time.RFC3339
		}
	}
}

func Init(opts ...Options) {
	c := &ZapConfig{}
	for _, opt := range opts {
		opt(c)
	}
	c = checkConfigOptions(c)
	timeEncode := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}
		if enc, ok := enc.(appendTimeEncoder); ok {
			enc.AppendTimeLayout(t, c.timeEncoderLayout)
			return
		}

		enc.AppendString(t.Format(c.timeEncoderLayout))
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        c.timeKey,
		LevelKey:       c.levelKey,
		NameKey:        c.nameKey,
		CallerKey:      c.callerKey,
		MessageKey:     c.messageKey,
		StacktraceKey:  c.stackTraceKey,
		LineEnding:     c.lineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeTime:     timeEncode,
	}
	config := &zap.Config{
		Level:         getLogLevel(c.level),
		Development:   c.development,
		Encoding:      c.encoding,
		EncoderConfig: encoderConfig,
		InitialFields: c.initialFields,
		OutputPaths:   c.outputPaths,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	switch config.Encoding {
	case ConsoleEncoder:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case JSONEncoder:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var ops []zap.Option
	stackLevel := zap.ErrorLevel
	if config.Development {
		stackLevel = zap.WarnLevel
	}
	if !config.DisableStacktrace {
		ops = append(ops, zap.AddStacktrace(stackLevel))
	}
	if len(config.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(config.InitialFields))
		keys := make([]string, 0, len(config.InitialFields))
		for k := range config.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, config.InitialFields[k]))
		}
		ops = append(ops, zap.Fields(fs...))
	}

	core := zapcore.NewCore(encoder, getLogWriter(config.OutputPaths), config.Level)
	logger := zap.New(core, ops...)
	logger.WithOptions()
	zap.ReplaceGlobals(logger)
}

func checkConfigOptions(config *ZapConfig) *ZapConfig {
	if config.encoding == "" {
		config.encoding = "json"
	}
	if config.level == "" {
		config.level = "debug"
	}
	if config.timeKey == "" {
		config.timeKey = "time"
	}
	if config.levelKey == "" {
		config.levelKey = "level"
	}
	if config.nameKey == "" {
		config.nameKey = "name"
	}
	if config.callerKey == "" {
		config.callerKey = "caller"
	}
	if config.messageKey == "" {
		config.messageKey = "msg"
	}
	if config.stackTraceKey == "" {
		config.stackTraceKey = "stack"
	}
	if config.lineEnding == "" {
		config.lineEnding = zapcore.DefaultLineEnding
	}
	if config.outputPaths == nil || len(config.outputPaths) == 0 {
		config.outputPaths = []string{"stdout"}
	}
	if config.callSkip == 0 {
		config.callSkip = 1
	}
	if config.timeEncoderLayout == "" {
		config.timeEncoderLayout = time.RFC3339
	}
	return config
}

func getLogLevel(level string) zap.AtomicLevel {
	level = strings.ToLower(level)
	var ret zap.AtomicLevel
	switch level {
	case "prod", "info":
		ret = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "dev", "debug":
		ret = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "test":
		ret = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		ret = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return ret
}

func getLogWriter(output []string) zapcore.WriteSyncer {
	writers := make([]zapcore.WriteSyncer, 0)
	for _, s := range output {
		var sync zapcore.WriteSyncer
		switch s {
		case "stdout":
			sync = zapcore.AddSync(os.Stdout)
		case "stderr":
			sync = zapcore.AddSync(os.Stderr)
		default:
			lumberJackLogger := &lumberjack.Logger{
				Filename:   s,
				MaxSize:    10,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   false,
			}
			sync = zapcore.AddSync(lumberJackLogger)
		}
		writers = append(writers, sync)
	}

	return zapcore.AddSync(zap.CombineWriteSyncers(writers...))
}

func Sync() {
	_ = zap.S().Sync()
	_ = zap.L().Sync()
}

func SPanic(args ...interface{}) {
	zap.S().Panic(args)
}

func SPanicf(msg string, args ...interface{}) {
	zap.S().Panicf(msg, args...)
}

func SFatal(args ...interface{}) {
	zap.S().Fatal(args)
}

func SFatalf(tpl string, args ...interface{}) {
	zap.S().Fatalf(tpl, args)
}

func SError(args ...interface{}) {
	zap.S().Error(args...)
}

func SErrorln(args ...interface{}) {
	zap.S().Errorln(args...)
}

func SErrorf(msg string, args ...interface{}) {
	zap.S().Errorf(msg, args...)
}

func SErrorw(msg string, keysAndValues ...interface{}) {
	zap.S().Errorw(msg, keysAndValues...)
}

func SInfo(args ...interface{}) {
	zap.S().Info(args...)
}

func SInfof(tpl string, args ...interface{}) {
	zap.S().Infof(tpl, args...)
}

func SInfoln(args ...interface{}) {
	zap.S().Infoln(args...)
}

func SDebug(args ...interface{}) {
	zap.S().Debug(args)
}

func SDebugln(args ...interface{}) {
	zap.S().Debugln(args)
}

func SDebugf(tpl string, args ...interface{}) {
	zap.S().Debugf(tpl, args...)
}

func SWarnf(tpl string, args ...interface{}) {
	zap.S().Warnf(tpl, args...)
}

func LInfo(msg string, fields ...zap.Field) {
	zap.L().Info(msg, fields...)
}

func LDebug(msg string, fields ...zap.Field) {
	zap.L().Debug(msg, fields...)
}
