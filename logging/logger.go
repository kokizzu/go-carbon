package logging

import (
	"bytes"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/Sirupsen/logrus"
)

var badMessageLogger *logrus.Logger

func init() {
	logrus.SetFormatter(&TextFormatter{})
	badMessageLogger = logrus.StandardLogger()
}

// SetLevel for default logger
func SetLevel(logger *logrus.Logger, lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	if logger == nil {
		logrus.SetLevel(level)
	} else {
		logger.Level = level
	}
	return nil
}

// PrepareFile creates logfile and set it writable for user
func PrepareFile(filename string, owner *user.User) error {
	if filename == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if fd != nil {
		fd.Close()
	}
	if err != nil {
		return err
	}
	if err := os.Chmod(filename, 0644); err != nil {
		return err
	}
	if owner != nil {

		uid, err := strconv.ParseInt(owner.Uid, 10, 0)
		if err != nil {
			return err
		}

		gid, err := strconv.ParseInt(owner.Gid, 10, 0)
		if err != nil {
			return err
		}

		if err := os.Chown(filename, int(uid), int(gid)); err != nil {
			return err
		}
	}

	return nil
}

// Test run callable with changed logging output
func Test(callable func(*bytes.Buffer)) {
	buf := &bytes.Buffer{}
	logrus.SetOutput(buf)

	callable(buf)

	logrus.SetOutput(os.Stderr)
}

// TestWithLevel run callable with changed logging output and log level
func TestWithLevel(level string, callable func(*bytes.Buffer)) {
	originalLevel := logrus.GetLevel()
	defer logrus.SetLevel(originalLevel)
	SetLevel(nil, level)

	Test(callable)
}
