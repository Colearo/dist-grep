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
	// Open config file
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

	// Create a new test log file
	file, err := os.OpenFile("./test_"+config.LocalInfo.NodeName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	logPrefix := "[" + config.LocalInfo.NodeName + "]: "
	logger := log.New(file, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)

	i := 0

	// Source for genearting random number
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	for {
		// Generating 1000 known pattern log
		if i < 1000 {
			logger.Printf("Gen logs here\n")
			logger.Printf("KNOWN PATTERN 1000 [%d]\n", i)
		} else {
			// Else generating 1000 random log
			logger.Printf("Gen logs here %d\n", randGen.Intn(10000))
			logger.Printf("RANDOM PATTERN [%d]\n", i)
			// Each 10 seconds do the log gen
			time.Sleep(10 * time.Second)
		}
		i += 1
	}
}
