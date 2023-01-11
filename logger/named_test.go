package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var byteBuffer = &bytes.Buffer{}

func TestWithName(t *testing.T) {
	assert := assert.New(t)
	logger := WithName("foo")
	assert.NotNil(logger)
	logger.Info("foo")
	assert.False(logger.Level().Enabled(zap.PanicLevel))
}

func TestHydrate(t *testing.T) {
	assert := assert.New(t)
	defer byteBuffer.Reset()
	logger := WithName("foo")
	assert.NotNil(logger)
	logger.Info("foo")
	assert.False(logger.Level().Enabled(zap.PanicLevel))
	Hydrate(newTestLogger())
	logger.Info("foo")
	assert.NoError(logger.Sync())
	assert.True(logger.Level().Enabled(zap.DebugLevel))
	assert.Equal(byteBuffer.String(), `{"level":"info","logger":"foo","msg":"foo"}`+"\n")
}

func newTestLogger(options ...zap.Option) *zap.Logger {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(byteBuffer), zap.DebugLevel)
	return zap.New(core).WithOptions(options...)
}
