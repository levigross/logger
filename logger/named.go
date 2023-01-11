package logger

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"go.uber.org/zap"
)

var (
	properLogger *zap.Logger
	hydrateOnce  sync.Once
	loggers      = make(map[string]*zap.Logger)
	loggersMutex sync.Mutex
)

// WithName should be called when creating top of page loggers
func WithName(name string) (l *zap.Logger) {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()
	if properLogger == nil {
		properLogger = zap.NewNop()
	}
	l = properLogger.Named(name)
	loggers[name] = l
	return l
}

func Hydrate(log *zap.Logger) {
	hydrateOnce.Do(func() {
		loggersMutex.Lock()
		defer loggersMutex.Unlock()
		// overwrite it because this is only called once
		atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&properLogger)), unsafe.Pointer(log))
		for name := range loggers {
			oldLogger := loggers[name]
			newLogger := properLogger.Named(name)
			*oldLogger = *newLogger
		}
	})
}
