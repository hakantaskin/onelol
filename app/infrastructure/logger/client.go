package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// New ...
func New(out io.Writer, env, logLevel string) *logrus.Logger {
	log := logrus.New()
	log.Out = out
	if env == "dev" {
		log.Formatter = &logrus.TextFormatter{}
	} else {
		log.Formatter = getJSONFormatter()
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.Level = level

	return log
}

// JSONFormatter Wrapper for logrus.JSONFormatter
type JSONFormatter struct {
	logrus.JSONFormatter
}

// getJSONFormatter
func getJSONFormatter() *JSONFormatter {
	jsonFormatter := logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05-0700",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
	}
	return &JSONFormatter{jsonFormatter}
}
