package yago

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type AppConfig struct {
	*viper.Viper
}

func IsEnvProd() bool {
	return Config.GetString("app.env") == "prod"
}

func IsEnvDev() bool {
	return Config.GetString("app.env") == "dev"
}

var Config *AppConfig

func NewAppConfig(cfgPath string) *AppConfig {
	cfg := &AppConfig{viper.New()}
	cfg.SetConfigFile(cfgPath)
	err := cfg.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n, \"--help\" gives usage information", err))
	}

	cfg.SetDefault("app.app_name", "APP")

	appName := cfg.GetString("app.app_name")
	appName = strings.ReplaceAll(appName, "-", "_")
	cfg.SetEnvPrefix(appName)
	cfg.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	cfg.SetEnvKeyReplacer(replacer)
	return cfg
}

func Hostname() string {
	endpoint := Config.GetString("endpoint")
	if endpoint == "" {
		endpoint, _ = os.Hostname()
	}
	return endpoint
}
