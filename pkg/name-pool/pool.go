package name_pool

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	index  uint64
	length uint64

	mutex sync.Mutex

	namePoolQueue []string
)

const (
	defaultQueueLength        = 5000
	defaultReplica            = 1
	defaultNamePoolConfigPath = "/configs/name-pool"
)

func Execute() error {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return err
	}
	configPath := basePath + defaultNamePoolConfigPath
	readConfig("config", configPath, "yaml")
	if configMessage.scope < 0 {
		namePoolQueue = make([]string, defaultQueueLength*defaultReplica)
		length = defaultQueueLength * defaultReplica
	} else {
		namePoolQueue = make([]string, configMessage.scope*configMessage.replica)
		length = configMessage.scope * configMessage.replica
	}
	// 初始化队列
	index = 0
	for i := uint64(1); i <= configMessage.replica; i++ {
		for j := uint64(1); j <= configMessage.scope; j++ {
			namePoolQueue[index] = configMessage.prefixName + strconv.FormatUint(j, 10)
			index++
		}
	}
	index = 0
	http.HandleFunc(configMessage.pattern, getName)
	if err := http.ListenAndServe(":"+configMessage.port, nil); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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
	fmt.Println(str)
	return str
}
