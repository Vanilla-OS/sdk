package logs

/*	License: GPLv3
	Authors:
		Mirko Brombin <brombin94@gmail.com>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2026
	Description: Vanilla OS SDK component.
*/

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

// Info logs an informational message to both console and file.
func (l *Logger) Info(msg string) {
	l.File.Info().Msg(msg)
	l.Term.Info().Msg(msg)
}

// Infof logs a formatted informational message to both console and file.
func (l *Logger) Infof(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.File.Info().Msg(msg)
	l.Term.Info().Msg(msg)
}

// Warn logs a warning message to both console and file.
func (l *Logger) Warn(msg string) {
	l.File.Warn().Msg(msg)
	l.Term.Warn().Msg(msg)
}

// Warnf logs a formatted warning message to both console and file.
func (l *Logger) Warnf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.File.Warn().Msg(msg)
	l.Term.Warn().Msg(msg)
}

// Error logs an error message to both console and file.
func (l *Logger) Error(msg string) {
	l.File.Error().Msg(msg)
	l.Term.Error().Msg(msg)
}

// Errorf logs a formatted error message to both console and file.
func (l *Logger) Errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.File.Error().Msg(msg)
	l.Term.Error().Msg(msg)
}

// Debug logs a debug message to both console and file.
func (l *Logger) Debug(msg string) {
	l.File.Debug().Msg(msg)
	l.Term.Debug().Msg(msg)
}

// Debugf logs a formatted debug message to both console and file.
func (l *Logger) Debugf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.File.Debug().Msg(msg)
	l.Term.Debug().Msg(msg)
}

// Trace logs a trace message to both console and file.
func (l *Logger) Trace(msg string) {
	l.File.Trace().Msg(msg)
	l.Term.Trace().Msg(msg)
}

// Tracef logs a formatted trace message to both console and file.
func (l *Logger) Tracef(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.File.Trace().Msg(msg)
	l.Term.Trace().Msg(msg)
}
