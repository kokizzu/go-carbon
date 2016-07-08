package receiver

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

var MessageLog = logrus.StandardLogger().Infof

func SetupMessageLog(logger *logrus.Logger, level string) error {
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
		return fmt.Errorf("unknown level %#v", level)
	}

	return nil
}
