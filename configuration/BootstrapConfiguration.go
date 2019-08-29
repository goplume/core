package configuration

import (
	"fmt"
	"github.com/goplume/core/utils/logger"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	"os"
	"regexp"
)

var LoadConfigurationError = ""

func BootstrapConfiguration() {

	LOGGER_STRATEGY := os.Getenv("LOGGER_STRATEGY")
	switch LOGGER_STRATEGY {
	case logger.CONSOLE_STRATEGY:
		logger.SetupLogrusToConsole(log.StandardLogger())
	default:
		logger.SetupLogrusToFile(log.StandardLogger(), "")
	}

	//SwaggerConfiguration()
	//gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	//gin.DisableConsoleColor()

	dir, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}
	log.Info("WORKDIR: " + dir)

	BootstrapConfiguration_Viper()
}

func BootstrapConfiguration_Viper() {

	config.SetConfigName("config")                                                // name of config file (without extension)
	config.AddConfigPath("/etc/mvisa_merchant")                                   // path to look for the config file in
	config.AddConfigPath("$HOME/.merchantapi")                                    // call multiple times to add many search paths
	config.AddConfigPath("$HOME/.mvisa_merchant")                                 // call multiple times to add many search paths
	config.AddConfigPath(".")                                                     // optionally look for config in the working directory
	config.AddConfigPath("..")                                                    // optionally look for config in the working directory
	config.AddConfigPath("config")                                                // optionally look for config in the working directory
	config.AddConfigPath("../config")                                             // optionally look for config in the working directory
	config.AddConfigPath("../../config")                                          // optionally look for config in the working directory
	config.AddConfigPath("../../../config")                                       // optionally look for config in the working directory
	config.AddConfigPath("../../../../config")                                    // optionally look for config in the working directory
	LoadConfigurationError := config.ReadInConfig()                               // Find and read the config file
	if LoadConfigurationError != nil {
		// Handle errors reading the config file
		log.Error(fmt.Errorf("Fatal error config file: %s \n", LoadConfigurationError))
		panic(LoadConfigurationError)
	}

	configuratioVersion := config.GetString("version")
	log.Info("Configuration version: " + configuratioVersion)

	//    log.Info(config.GetViper().AllKeys())
	//    log.Info(config.GetViper().AllSettings())
}

// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`([a-z_0-9])([A-Z])`)

//type conventionalMarshaller struct {
//    Value interface{}
//}
//
//func (self conventionalMarshaller) MarshalJSON() ([]byte, error) {
//    marshalled, err := json.Marshal(self.Value)
//
//    converted := keyMatchRegex.ReplaceAllFunc(
//        marshalled,
//        func(match []byte) []byte {
//            return bytes.ToLower(wordBarrierRegex.ReplaceAll(
//                match,
//                []byte(`${1}_${2}`),
//            ))
//        },
//    )
//
//    return converted, err
//}
