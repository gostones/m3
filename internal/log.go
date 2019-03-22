package internal

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	//logrus.SetFormatter(&logrusrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(logFile())

	//Logrus has six logging levels: Debug, Info, Warning, Error, Fatal and Panic.
	//
	level := logLevel()
	logrus.SetLevel(level)

	//
	const app_name = "m3"
	logger = logrus.WithFields(logrus.Fields{
		"application_name": app_name,
	})

	//
	logger.Infof("Logrus initialized. log level: %s", level)
}

//default to debug if env not set
func logLevel() (level logrus.Level) {
	l := os.Getenv("log_level")
	if l == "" {
		l = "debug"
	}
	level, err := logrus.ParseLevel(l)
	if err != nil {
		level = logrus.DebugLevel
	}
	return
}

func logFile() io.Writer {
	f := os.Getenv("log_file")
	if f == "" {
		return os.Stdout
	}
	ensureDir(f)
	w, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		return w
	} else {
		return os.Stdout
	}
}

func ensureDir(name string) {
	dir := filepath.Dir(name)
	if dir == "" {
		return
	}
	if _, err := os.Stat(dir); err == nil {
		return
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
}

// Logger exports logger
func Logger() *logrus.Entry {
	return logger
}
