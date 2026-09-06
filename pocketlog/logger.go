package pocketlog

type Logger struct{}

func (l *Logger) DebugF(format string, args ...any) {}

func (l *Logger) Infof(format string, args ...any) {
}
