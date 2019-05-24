package yago

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type AppConfig struct {
	*viper.Viper
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

func init() {
	defaultDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cfgPath := flag.String("c", fmt.Sprintf("%s/app.toml", defaultDir), "config file path")
	flag.Parse()
	Config = NewAppConfig(*cfgPath)
}
