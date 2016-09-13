package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils"
	"github.com/ysqi/atop/common/log2"
	"github.com/ysqi/atop/common/util"
)

var (
	// AppCfg is the instance of Config, store the config information from file
	AppCfg *myAppConfig
	// AppPath is the absolute path to the app
	AppPath string
	// AppCfgPath is the path to the config files
	AppCfgPath string
)

func init() {
	var err error
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	if p := os.Getenv("AppConfigPath"); p != "" && utils.FileExists(p) {
		AppCfgPath = p
	} else {
		workPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		AppCfgPath = filepath.Join(workPath, "conf", "app.conf")
	}
	if !utils.FileExists(AppCfgPath) {
		AppCfgPath = filepath.Join(AppPath, "conf", "app.conf")
		if !utils.FileExists(AppCfgPath) {
			if path, err := util.SearchFile(filepath.Join("conf", "app.conf")); err == nil {
				AppCfgPath = path
			} else {
				AppCfgPath = ""
				AppCfg = &myAppConfig{
					innerConfig: config.NewFakeConfig(),
					RunMode:     "dev",
				}
				return
			}
		}
	}
	log2.Infof("App Config Path=%q", AppCfgPath)
	if err = parseConfig(AppCfgPath); err != nil {
		panic(err)
	}
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	AppCfg, err = newAppConfig("ini", appConfigPath)
	if err != nil {
		return err
	}
	return assignConfig(AppCfg)
}

func assignConfig(ac config.Configer) error {
	// set the run mode first
	if envRunMode := os.Getenv(ac.DefaultString("AppName", "MyAPP") + "_RunMode"); envRunMode != "" {
		AppCfg.RunMode = envRunMode
	} else if runMode := ac.String("RunMode"); runMode != "" {
		AppCfg.RunMode = runMode
	}
	AppCfg.AppName = ac.String("AppName")

	logOuputs := make(map[string]string)
	if lo := ac.String("LogOutputs"); lo != "" {
		los := strings.Split(lo, ";")
		for _, v := range los {
			if logType2Config := strings.SplitN(v, ",", 2); len(logType2Config) == 2 {
				logOuputs[logType2Config[0]] = logType2Config[1]
			} else {
				continue
			}
		}
	}

	//init log
	logs.Reset()
	for adaptor, config := range logOuputs {
		err := logs.SetLogger(adaptor, config)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("%s with the config %q got err:%s", adaptor, config, err.Error()))
		}
	}
	logs.SetLogFuncCall(true)
	return nil
}

type myAppConfig struct {
	innerConfig config.Configer
	RunMode     string
	AppName     string
}

func newAppConfig(appConfigProvider, appConfigPath string) (*myAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &myAppConfig{innerConfig: ac, RunMode: "dev"}, nil
}

func (b *myAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(b.RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *myAppConfig) String(key string) string {
	if v := b.innerConfig.String(b.RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *myAppConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(b.RunMode + "::" + key); len(v) > 0 {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *myAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(b.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *myAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(b.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *myAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(b.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *myAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(b.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *myAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *myAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *myAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *myAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
