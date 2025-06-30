package logs

import (
	"fmt"
	"sync"

	"github.com/phuslu/log"
)

// Logger represents a logger for the application.
type Logger struct {
	// Term is the logger used to log messages to the console, use this for any
	// logging the user should see.
	Term log.Logger

	// File is the logger used to log messages to the vlogs directory, use this
	// for any internal logging the user doesn't need to see.
	File log.Logger

	mu       sync.Mutex
	ErrIndex map[string]int
}

func (l *Logger) nextErrIndex(prefix string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ErrIndex == nil {
		l.ErrIndex = make(map[string]int)
	}

	idx := l.ErrIndex[prefix]
	l.ErrIndex[prefix] = idx + 1
	return idx
}

// InfoCtx logs an informational message using the provided context.
func (l *Logger) InfoCtx(ctx *LogContext, msg string) {
	prefix := ctx.Prefix()
	formatted := fmt.Sprintf("%s:info:%s", prefix, msg)
	l.File.Info().Msg(formatted)
	l.Term.Info().Msg(formatted)
}

// WarnCtx logs a warning message using the provided context.
func (l *Logger) WarnCtx(ctx *LogContext, msg string) {
	prefix := ctx.Prefix()
	formatted := fmt.Sprintf("%s:warn:%s", prefix, msg)
	l.File.Warn().Msg(formatted)
	l.Term.Warn().Msg(formatted)
}

// ErrorCtx logs an error message using the provided context. The error index
// is automatically incremented per-context to provide consistent progression.
func (l *Logger) ErrorCtx(ctx *LogContext, msg string) {
	prefix := ctx.Prefix()
	idx := l.nextErrIndex(prefix)
	formatted := fmt.Sprintf("%s:err(%d):%s", prefix, idx, msg)
	l.File.Error().Msg(formatted)
	l.Term.Error().Msg(formatted)
}
