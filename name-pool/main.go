package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type config struct {
	prefixName string
	scope      int
	replica    int
}

var once sync.Once

var configMessage *config

var namePoolQueue []string

var (
	index  int
	length int
)

var mutex sync.Mutex

const (
	defaultQueueLength = 5000
	defaultReplica     = 1
)

func readConfig(configName string, configPath string, configType string) {
	once.Do(func() {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName(configName)
		viper.AddConfigPath(path + configPath)
		viper.SetConfigType(configType)
		err = viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		configMessage = &config{
			prefixName: viper.GetString("config.prefixName"),
			scope:      viper.GetInt("config.scope"),
			replica:    viper.GetInt("config.replica"),
		}
	})
}

func main() {
	readConfig("config", "/name-pool/conf", "yaml")
	if configMessage.scope < 0 {
		namePoolQueue = make([]string, defaultQueueLength*defaultReplica)
		length = defaultQueueLength * defaultReplica
	} else {
		namePoolQueue = make([]string, configMessage.scope*configMessage.replica)
		length = configMessage.scope * configMessage.replica
	}
	index = 0
	for i := 1; i <= configMessage.replica; i++ {
		for j := 1; j <= configMessage.scope; j++ {
			namePoolQueue[index] = configMessage.prefixName + strconv.Itoa(j)
			index++
		}
	}
	index = 0
	http.HandleFunc("/name", getName)
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}

}

func getName(w http.ResponseWriter, r *http.Request) {
	str := readName()
	w.Write([]byte(str))
}

func readName() string {
	str := ""
	mutex.Lock()
	if index < length {
		str = namePoolQueue[index]
		index++
	}
	mutex.Unlock()
	return str
}
