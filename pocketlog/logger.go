package pocketlog

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	threshold Level
	output    io.Writer
}

// Debug formats and prints a message if the log level is debug or higher
func (l *Logger) Debugf(format string, args ...any) {
	if l.output == nil {
		l.output = os.Stdout
	}
	logf(format, l, LevelDebug, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	logf(format, l, LevelInfo, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	logf(format, l, LevelError, args...)
}

func logf(format string, l *Logger, lvl Level, args ...any) {
	if l.threshold > lvl {
		return
	}
	_, _ = fmt.Fprintf(l.output, format, args...)
}

// New returns you a logger, ready to log at the required threshold.
func New(threshold Level, output io.Writer) *Logger {
	return &Logger{
		threshold: threshold,
		output:    output,
	}
}
