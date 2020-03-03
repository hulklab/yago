package yago

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
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

func GetAppName() string {
	return Config.GetString("app.app_name")
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
	//cfg.WatchConfig()
	//cfg.OnConfigChange(func(e fsnotify.Event) {
	//	// viper配置发生变化了 执行响应的操作
	//	fmt.Println("Config file changed:", e.Name, e.String())
	//})

	return cfg
}

func Hostname() string {
	endpoint := Config.GetString("endpoint")
	if endpoint == "" {
		endpoint, _ = os.Hostname()
	}
	return endpoint
}

func defaultCfgPath() string {
	defaultDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return fmt.Sprintf("%s/app.toml", defaultDir)
}

var cfgLock = new(sync.Mutex)

func initConfig() {
	cfgPath := flag.String("c", defaultCfgPath(), "config file path")
	_ = flag.Bool("h", false, "help")
	_ = flag.Bool("help", false, "help")
	flag.Parse()
	Config = NewAppConfig(*cfgPath)
}

func ReloadConfig() error {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	// 重新加载配置文件
	err := Config.ReadInConfig()
	if err != nil {
		return err
	}

	// 清理组件
	Component.clear()
	return nil
}

func getPidFile() (string, bool) {
	pidfile := Config.GetString("app.pidfile")
	if pidfile == "" {
		return "", false
	}
	return pidfile, true
}
