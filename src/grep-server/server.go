package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Config struct {
	LocalInfo LocalInfo
}
type LocalInfo struct {
	NodeName   string
	ServerPort string
	LogPath    string
}

var config Config

func main() {
	configFile, err := os.Open("../../config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	// Read json file's contents and cache them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(configBytes, &config)

	serverAddr := ":" + config.LocalInfo.ServerPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	printError(err)
	listen, err := net.ListenTCP("tcp", tcpAddr)
	printError(err)
	for {
		connect, err := listen.Accept()
		if err != nil {
			continue
		}
		handleMsg(connect)
	}
}

func handleMsg(connect net.Conn) {
	buffer := make([]byte, 1024)
	n, err := connect.Read(buffer)
	data := string(buffer[:n])
	fmt.Println(data)
	commands := strings.Fields(data)
	if err != nil {
		fmt.Println("Fatal Error")
		return
	}
	validFlagC := regexp.MustCompile(`^\-[a-zA-Z]*c`)
	flagC := false
	for _, val := range commands {
		if validFlagC.MatchString(val) {
			flagC = true
			break
		}
	}

	commands = append(commands, "-Hn", config.LocalInfo.LogPath)
	cmd := exec.Command("grep", commands...)
	stdOut, err := cmd.StdoutPipe()
	var stdOutErr []byte
	var wg sync.WaitGroup
	if flagC == false {
		wg.Add(1)
		go func() {
			commandsCount := append(commands, "-c")
			cmdCount := exec.Command("grep", commandsCount...)
			stdOutErr, _ = cmdCount.CombinedOutput()
			fmt.Println(string([]byte(stdOutErr)))
			wg.Done()
		}()
	}
	printError(err)
	//stdOutErr, _ := cmd.CombinedOutput()
	bufferOut := make([]byte, 1024)
	err = cmd.Start()
	printError(err)
	for {
		lenOut, errOut := stdOut.Read(bufferOut)

		if errOut != nil {
			if errOut == io.EOF {
				break
			} else {
				printError(errOut)
			}
		}

		fmt.Printf("%s", string(bufferOut[:lenOut]))
		connect.Write(bufferOut[:lenOut])
	}

	err = cmd.Wait()
	printError(err)

	if flagC == false {
		wg.Wait()
		connect.Write([]byte(stdOutErr))
	}

	fmt.Println("[Debug] Output Ended")
	connect.Close()
}

func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]", err.Error())
	}
}
