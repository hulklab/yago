package yago

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hulklab/yago/libs/arr"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type AppConfig struct {
	*viper.Viper
}

func (c *AppConfig) ReadFileConfig(cfgPath string) error {
	err := c.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: %s \n, \"--help\" gives usage information", err)
	}

	// deal with import file
	importFiles, err := c.readImportFiles(cfgPath)
	if err != nil {
		return fmt.Errorf("Fatal error merge import config file: %s ", err)
	}

	if len(importFiles) >= 2 {
		log.Println("import configs:", importFiles)
		// the last one don't need merge
		for i := len(importFiles) - 2; i >= 0; i-- {
			importFile := importFiles[i]

			c.SetConfigFile(importFile)
			_ = c.MergeInConfig()
		}
	}
	return err
}

func (c *AppConfig) readImportFiles(cfgPath string) ([]string, error) {
	if !c.IsSet("import") {
		return nil, nil
	}

	importFiles := make([]string, 0)
	// put current file into the head of list
	importFiles = append(importFiles, cfgPath)

	for {
		includeFile := c.GetString("import")
		if !filepath.IsAbs(includeFile) {
			includeFile, _ = filepath.Abs(filepath.Join(filepath.Dir(cfgPath), includeFile))
		}

		if arr.InArray(includeFile, importFiles) {
			return importFiles, fmt.Errorf("circle import config file")
		}

		importFiles = append(importFiles, includeFile)

		c.SetConfigFile(includeFile)

		err := c.ReadInConfig()
		if err != nil {
			return importFiles, fmt.Errorf("Fatal error merge include config file: %s ", err)
		}

		if !c.IsSet("import") {
			break
		}
	}

	return importFiles, nil
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

	if strings.HasPrefix(cfgPath, "etcd+") {
		// only support etcd now
		cs := strings.Split(cfgPath, "+")
		if len(cs) != 2 {
			panic("remote config only support like this: etcd+http://127.0.0.1:2379/yago/conf/app.toml")
		}

		u, err := url.Parse(cs[1])
		if err != nil {
			panic("parse config err: " + err.Error())
		}

		err = cfg.AddRemoteProvider("etcd", u.Scheme+"://"+u.Host, u.Path)
		if err != nil {
			panic(err)
		}

		// 获取 path 后缀
		ext := filepath.Ext(u.Path)
		if len(ext) > 1 {
			cfg.SetConfigType(ext[1:])
		} else {
			cfg.SetConfigType("toml")
		}

		err = cfg.ReadRemoteConfig()
		if err != nil {
			panic(err)
		}

		err = cfg.WatchRemoteConfigOnChannel()
		if err != nil {
			log.Println("watch remote config err:", err.Error())
		}
	} else {
		cfg.SetConfigFile(cfgPath)
		err := cfg.ReadFileConfig(cfgPath)
		if err != nil {
			panic(err)
		}
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

func defaultCfgPath() string {
	defaultDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return fmt.Sprintf("%s/app.toml", defaultDir)
}

var cfgPath *string
var cfgLock = new(sync.Mutex)

func initConfig() {
	cfgPath = flag.String("c", defaultCfgPath(), "config file path")
	_ = flag.Bool("h", false, "help")
	_ = flag.Bool("help", false, "help")
	flag.Parse()
	Config = NewAppConfig(*cfgPath)
}

func reloadConfig() error {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	// reload file
	err := Config.ReadFileConfig(*cfgPath)
	if err != nil {
		return err
	}

	// clear components config
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
