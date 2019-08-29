package logger

import (
	"github.com/goplume/core/fault"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"runtime"
)

var (
	// DefaultLogLevels is the log levels for which errors are reported by Hook, if Hook.LogLevels is not set.
	DefaultLogLevels = []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
)

func NewELKLogrusHook(Url string) ELKLogrusHook {
	client := resty.New()
	client.SetDebug(false)
	//client.SetDebug(true)
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "mvisa-merchant", //
	})

	return ELKLogrusHook{
		Url:        Url,
		HttpClient: client,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				//logrus.FieldKeyTime:  "@timestamp",
				//logrus.FieldKeyLevel: "@level",
				//logrus.FieldKeyMsg:   "@message",
				//logrus.FieldKeyFunc:  "@caller",
				logrus.FieldKeyTime: "LogDate",
				logrus.FieldKeyFunc: "func",
				logrus.FieldKeyFile: "category",
				//logrus.FieldKeyLevel: "Level",
				//logrus.FieldKeyMsg:   "Message",
				//logrus.FieldKeyFunc:  "@caller",
			},
			CallerPrettyfier: func(caller *runtime.Frame) (function string, file string) {
				file, function = ParseFunctionName(caller.PC)
				return function, file
			},
		},
	}
}

type ELKLogrusHook struct {
	LogLevels  []logrus.Level
	HttpClient *resty.Client
	Url        string
	Formatter  *logrus.JSONFormatter
}

func (this ELKLogrusHook) Levels() []logrus.Level {
	levels := this.LogLevels
	if levels == nil {
		levels = DefaultLogLevels
	}
	return levels
}

// Fire reports the log entry as an error to the APM Server.
func (this ELKLogrusHook) Fire(entry *logrus.Entry) error {
	if entry.Caller != nil {
		//functionName := ParseFunctionName(entry.Caller.PC)
		//entry = entry.WithField("category", functionName)
	}
	request := this.HttpClient.R()
	//request.SetPathParams(pathParams)
	//body, e := json.Marshal(entry.Data)
	//if e != nil {
	//	return e
	//}
	bytes, e := this.Formatter.Format(entry)
	if e != nil {
		return e
	}
	request.SetBody(bytes)
	resp, err := request.Post(this.Url)
	if err != nil {
		return err
	}

	if resp != nil && resp.IsError() {
		return fault.IntegratrionExceptionClientErrorR(resp)
	}

	return nil
}
