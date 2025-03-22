package helpers

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func SetupLogging() {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(logrus.StandardLogger().Out)
	logrus.SetReportCaller(true)

	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevelStr = strings.ToLower(logLevelStr)

	logLevels := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}

	logLevel := logLevels[logLevelStr]

	logrus.SetLevel(logLevel)
}
