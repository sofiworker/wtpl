package conf

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
	"wtpl/pkg/logger"
	"wtpl/pkg/os"
)

var config = new(Config)

// InitConfig 配置默认顺序为：配置文件 -> 环境变量 -> 命令行参数
// 优先级为从右到左
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/")
	viper.AddConfigPath("/etc/程序名称/conf/")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config = GetDefaultConfig()
		} else {
			panic(err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(config, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "yaml"
	})
	if err != nil {
		panic(err)
	}
	if config.App.Name == "" {
		config.App.Name = os.MustGetAppName()
	}
	//logger.Init(logger.WithLogEncoding("console"))
	logger.SDebugf("%v", *config)
}

func GetDefaultConfig() *Config {
	return &Config{}
}

func GetConfig() *Config {
	return config
}
