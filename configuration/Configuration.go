package configuration

import (
    "fmt"
    "github.com/goplume/core/utils/logger"
    config "github.com/spf13/viper"
    "regexp"
)

type Configuration struct {
    ConfigurationContext string
    Log                  *logger.Logger
}

func NewServiceConfiguration(
    serviceName,
    ConfigurationContext string,
    Log *logger.Logger,
) Configuration {
    return Configuration{
        ConfigurationContext: "services." + serviceName + "." + ConfigurationContext,
        Log: Log,
    }
}

func NewGlobalConfiguration(ConfigurationContext string) Configuration {
    return Configuration{
        ConfigurationContext: ConfigurationContext,
    }
}

func (this Configuration) GetBool(key string) bool {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetBool(path)
    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %v: %v", path, value))
    }
    return value
}

func (this Configuration) ReadConfigValue(
    key string,
    rf func() string,
    df func() string,
) {
    if config.IsSet(key) {
        value := rf()
        //rlog.Info(fmt.Sprintf("Read %s: %s", key, value))
        if this.Log != nil && this.Log.RLog != nil {
            this.Log.RLog.Info(fmt.Sprintf("Read %v: %v", key, value))
        }
    } else {
        defaultValue := df()
        if this.Log != nil && this.Log.RLog != nil {
            this.Log.RLog.Info(fmt.Sprintf("Read %v: %v", key, defaultValue))
        }
    }

}

func (this Configuration) GetBoolD(key string, defaultValue bool) (value bool) {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    this.ReadConfigValue(path, func() string {
        value = config.GetBool(path)
        return fmt.Sprintf("%v", value)
    }, func() string {
        value = defaultValue
        return fmt.Sprintf("%v", defaultValue)
    })
    return value
}

func (this Configuration) IsSet(key string) bool {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.IsSet(path)
    //rlog.Info(fmt.Sprintf("IsSet ", path, ":", value))
    return value
}

func (this Configuration) GetInt(key string) int {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetInt(path)
    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %s: %d", path, value))
    }
    return value
}

func (this Configuration) GetInt32(key string) int32 {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetInt32(path)
    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %s: %d", path, value))
    }
    return value
}

func (this Configuration) GetUint64(key string) uint64 {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetUint64(path)
    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %s: %d", path, value))
    }
    return value
}

func (this Configuration) GetString(key string) string {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetString(path)
    loggedValue := value
    if this.IsSensitiveParameter(key) {
        loggedValue = this.MaskSensitiveValue(value)
    }

    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %s: %s", path, loggedValue))
    }

    return value
}

func (this Configuration) GetStrings(key string) []string {
    path := fmt.Sprintf(this.ConfigurationContext, key)
    value := config.GetStringSlice(path)
    loggedValue := value
    if this.IsSensitiveParameter(key) {
        loggedValue = []string{"****"}
    }

    if this.Log != nil && this.Log.RLog != nil {
        this.Log.RLog.Info(fmt.Sprintf("Read %s: %s", path, loggedValue))
    }

    return value
}

func (this Configuration) MaskSensitiveValue(parameterValue string) string {
    if len(parameterValue) == 0 {
        return "value is empty"
    }

    return "*****"
}

func (this Configuration) IsSensitiveParameter(parameterName string) bool {
    var likePasswd = regexp.MustCompile(`.*passw.*`)
    var likeSecret = regexp.MustCompile(`.*secret.*`)
    if len(parameterName) == 0 {
        return false
    }

    return likePasswd.MatchString(parameterName) || likeSecret.MatchString(parameterName)
}

func (this Configuration) GetMap(key string) map[string]interface{} {
    path := this.keyPath(key)
    ma := config.Get(path)
    if ma != nil {
        value := ma.(map[string]interface{})
        return value
    }
    return nil
}

func (this Configuration) keyPath(key string) string {
    return fmt.Sprintf(this.ConfigurationContext, key)
}
