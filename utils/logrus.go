package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

var log = logrus.New()

func InitElasticAPM() {
	logrus.AddHook(&apmlogrus.Hook{})
}

func InitLogger() {
	log.Out = os.Stdout

	log.SetLevel(logrus.InfoLevel)

	log.SetFormatter(&logrus.JSONFormatter{})

	log.AddHook(&apmlogrus.Hook{})
}

func GetLogger() *logrus.Logger {
	return log
}
