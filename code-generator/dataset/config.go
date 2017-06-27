package GoConfig

import (
	"fmt"

	"github.com/spf13/viper"
)

var allowedConfigTypes = map[string]bool{
	"json": true,
	"yaml": true,
	"toml": true,
}

type ConfigOptions struct {
	ConfigType string
	ConfigFile string
	ConfigPath string
}

func Init(options *ConfigOptions) error {
	if !isConfigTypeAllowed(options.ConfigType) {
		return fmt.Errorf("Config type of %s not allowed", options.ConfigType)
	}
	viper.WatchConfig()
	setConfigOptions(options)
	return viper.ReadInConfig()
}

func isConfigTypeAllowed(configType string) bool {
	if _, ok := allowedConfigTypes[configType]; !ok {
		return false
	}
	return true
}

func setConfigOptions(options *ConfigOptions) {
	viper.SetConfigType(options.ConfigType)
	viper.SetConfigName(options.ConfigFile)
	viper.AddConfigPath(options.ConfigPath)
}

func GetConfigStringValue(key string) string {
	return viper.GetString(key)
}

func GetConfigIntValue(key string) int {
	return viper.GetInt(key)
}

func GetConfigFloatValue(key string) float64 {
	return viper.GetFloat64(key)
}

func GetConfigBoolValue(key string) bool {
	return viper.GetBool(key)
}

func GetConfigMapValue(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func SetConfigValue(key string, value interface{}) {
	viper.Set(key, value)
}

func HasKey(key string) bool {
	return viper.IsSet(key)
}
