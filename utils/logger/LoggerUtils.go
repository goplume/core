package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
)

func SetupLogrusToFile(log *logrus.Logger, url string) {
	log.Info("Setup log to file")

	log.SetReportCaller(true)

	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
		NoColors:    true,
	})

	//log.SetFormatter(&logrus.TextFormatter{
	//   DisableColors: true,
	//   FullTimestamp: true,
	//})

	if len(url) > 0 {
		log.AddHook(NewELKLogrusHook(url))
	}

	// You could set this to any `io.Writer` such as a file
	logFile, err := os.OpenFile("service.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
	}

	//log.SetOutput(logFile)

	log.AddHook(lfshook.NewHook(
		logFile,
		&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
			NoColors:    true,
		},
	))

	//gin.DefaultWriter = logFile
	//gin.DefaultErrorWriter = logFile
	log.Info("Logger out redirect to file")
}

func
SetupLogrusToConsole(log *logrus.Logger) {
	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "category"},
		NoColors:    true,
	})

}
