package log

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestConfigure(t *testing.T) {
	testCases := []struct {
		name              string
		level             string
		format            string
		expectedLevel     logrus.Level
		expectedFormatter logrus.Formatter
	}{
		{
			name:              "debug level with json format",
			level:             "debug",
			format:            "json",
			expectedLevel:     logrus.DebugLevel,
			expectedFormatter: &logrus.JSONFormatter{},
		},
		{
			name:              "info level with text format",
			level:             "info",
			format:            "text",
			expectedLevel:     logrus.InfoLevel,
			expectedFormatter: &logrus.TextFormatter{},
		},
		{
			name:              "warn level with json format",
			level:             "warn",
			format:            "json",
			expectedLevel:     logrus.WarnLevel,
			expectedFormatter: &logrus.JSONFormatter{},
		},
		{
			name:              "error level with text format",
			level:             "error",
			format:            "text",
			expectedLevel:     logrus.ErrorLevel,
			expectedFormatter: &logrus.TextFormatter{},
		},
		{
			name:              "unknown level with text format",
			level:             "unknown",
			format:            "text",
			expectedLevel:     logrus.InfoLevel,
			expectedFormatter: &logrus.TextFormatter{},
		},
		{
			name:              "test config",
			level:             "debug",
			format:            "json",
			expectedLevel:     logrus.DebugLevel,
			expectedFormatter: &logrus.JSONFormatter{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Configure(tc.level, tc.format)

			if logrus.GetLevel() != tc.expectedLevel {
				t.Errorf("Expected log level %v, but got %v", tc.expectedLevel, logrus.GetLevel())
			}

			if _, ok := logrus.StandardLogger().Formatter.(*logrus.JSONFormatter); ok {
				if _, ok := tc.expectedFormatter.(*logrus.JSONFormatter); !ok {
					t.Errorf("Expected JSON formatter, but got text")
				}
			} else if _, ok := logrus.StandardLogger().Formatter.(*logrus.TextFormatter); ok {
				if _, ok := tc.expectedFormatter.(*logrus.TextFormatter); !ok {
					t.Errorf("Expected text formatter, but got JSON")
				}
			}
		})
	}
}
