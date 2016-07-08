package receiver

import "github.com/Sirupsen/logrus"

var MessageLog = logrus.StandardLogger().Infof

func SetupMessageLog(logger *logrus.Logger, level string) {
	switch level {
	case "debug":
		MessageLog = logger.Debugf
	case "info":
		MessageLog = logger.Infof
	case "warn", "warning":
		MessageLog = logger.Warnf
	case "error":
		MessageLog = logger.Errorf
	default:
		MessageLog = logger.Infof
	}
}
