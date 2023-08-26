package component

import (
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"os"
)

var Log = initializeLogger()

func initializeLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.Out = file

	log.SetFormatter(&ecslogrus.Formatter{})
	return log
}
