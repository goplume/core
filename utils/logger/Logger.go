package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
	"strings"
)

type Logger struct {
	Component  string
	RLog       logrus.FieldLogger
	Strategy   LoggingStrategy
	ElasticUrl string
	// elk,other
	//ElasticRestClient *rest_client.MerchantServiceRestClient
}

type LoggingStrategy string

const (
	ELK_STRATEGY      LoggingStrategy = "ELK_STRATEGY"
	FILE_ELK_STRATEGY                 = "FILE_ELK_STRATEGY"
	CONSOLE_STRATEGY                  = "CONSOLE_STRATEGY"
	FILE_STRATEGY                     = "FILE_STRATEGY"
)

func (this *Logger) Out() (out *io.Writer) {
	return &logrus.StandardLogger().Out
}

var trunc_func_preifx = ""
var trunc_func_preifx_len = len(trunc_func_preifx)

func formatMessage(message string, pc uintptr, file string, line int, ok bool) string {
	function := runtime.FuncForPC(pc)
	functionName := function.Name()

	if len(functionName) > trunc_func_preifx_len {
		functionName = functionName[trunc_func_preifx_len:]
	}
	return fmt.Sprintf("%[2]s:%[4]d  - %[1]s", message, functionName, file, line, ok)
}

func ParseFunctionName(
	pc uintptr,
) (categoryName string, functionName string) {
	caller := runtime.FuncForPC(pc)
	callerName := caller.Name()

	if len(callerName) > trunc_func_preifx_len {
		callerName = callerName[trunc_func_preifx_len:]
	}

	lastSlash := strings.LastIndex(callerName, "/")
	if lastSlash > 0 {
		functionName = callerName[lastSlash+1:]
		categoryName = callerName[:lastSlash]

	} else {
		functionName = callerName
		categoryName = ""

	}

	return
}
