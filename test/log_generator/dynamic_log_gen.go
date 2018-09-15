package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/user"
	"time"
)

// Config file parser structure
type Config struct {
	LocalInfo LocalInfo
}
type LocalInfo struct {
	NodeName    string
	ServerPort  string
	LogPath     string
	TestLogPath string
}

var config Config

func main() {
	usr, _ := user.Current()
	usrHome := usr.HomeDir
	configFile, err := os.Open(usrHome + "/go/src/dist-grep/config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[ERROR]", err.Error())
	}
	defer configFile.Close()

	// Read json file's contents and cache them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(configBytes, &config)

	file, err := os.OpenFile("./test_"+config.LocalInfo.NodeName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	logPrefix := "[" + config.LocalInfo.NodeName + "]: "
	logger := log.New(file, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	i := 0

	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	for {
		if i < 1000 {
			logger.Println("Gen logs here")
		} else {
			logger.Printf("Gen logs here %d\n", randGen.Intn(10000))
		}
		i += 1
		time.Sleep(10 * time.Second)
	}
}
