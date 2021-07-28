package name_pool

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	once          sync.Once
	configMessage *config
)

type config struct {
	prefixName string
	scope      uint64
	replica    uint64
	port       string
	pattern    string
}

func readConfig(configName string, configPath string, configType string) {
	once.Do(func() {
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)
		viper.SetConfigType(configType)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s \n", err))
		}
		configMessage = &config{
			prefixName: viper.GetString("config.prefixName"),
			scope:      viper.GetUint64("config.scope"),
			replica:    viper.GetUint64("config.replica"),
			port:       viper.GetString("config.port"),
			pattern:    viper.GetString("config.pattern"),
		}
	})
}
