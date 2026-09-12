package pocketlog_test

import (
	"go-pocket-sized-projects/pocketlog"
	"testing"
)

type testWriter struct {
	contents string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.contents += string(p)
	return len(p), nil
}

const (
	debugMessage = "Why write I still all one, ever the same,"
	infoMessage  = "And keep invention in a noted weed,"
	errorMessage = "That every word doth almost tell my name,"
)

func TestLogger_DebugfInfofErrorf(t *testing.T) {
	type testCase struct {
		level    pocketlog.Level
		expected string
	}

	tt := map[string]testCase{
		"debug": {
			level: pocketlog.LevelDebug,
			expected: debugMessage + "\n" +
				infoMessage + "\n" +
				errorMessage + "\n",
		},
		"info": {
			level:    pocketlog.LevelInfo,
			expected: infoMessage + "\n" + errorMessage + "\n",
		},
		"error": {
			level:    pocketlog.LevelError,
			expected: errorMessage + "\n",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tw := &testWriter{}
			testLogger := pocketlog.New(tc.level, pocketlog.WithOutput(tw))
			testLogger.Debugf(debugMessage)
			testLogger.Infof(infoMessage)
			testLogger.Errorf(errorMessage)
			if tw.contents != tc.expected {
				t.Errorf("invalid contents, expected %q,got %q", tc.expected, tw.contents)
			}
		})
	}

}
