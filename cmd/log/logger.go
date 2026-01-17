package log

import (
	"github.com/sirupsen/logrus"
)

var levelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func Configure(level, format string) {
	logLevel, ok := levelMap[level]
	if !ok {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
}
