package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"sync"
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
	printError(err)
	defer configFile.Close()

	// Read json file's contents and cache them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(configBytes, &config)

	// Bind server address and port
	serverAddr := ":" + config.LocalInfo.ServerPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	printError(err)
	// Listen the request from client
	listen, err := net.ListenTCP("tcp", tcpAddr)
	printError(err)
	for {
		// Accept a TCP request
		connect, err := listen.Accept()
		if err != nil {
			continue
		}
		// Handle the message sent from client in coroutine
		go handleMsg(connect)
	}
}

func handleMsg(connect net.Conn) {
	// Making a buffer to accept the grep command content from client
	buffer := make([]byte, 1024)
	n, err := connect.Read(buffer)
	// Error means cannot read content to buffer
	if err != nil {
		fmt.Println("Fatal Error")
		return
	}
	// Convert to string
	data := string(buffer[:n])
	fmt.Println(data)
	// Split into array of arguments of 'grep' command
	commands := strings.Fields(data)

	// If command argument contains -c,
	// we do not need to try to count lines self
	validFlagC := regexp.MustCompile(`^\-[a-zA-Z]*c`)
	flagC := false
	for _, val := range commands {
		if validFlagC.MatchString(val) {
			flagC = true
			break
		}
	}

	// Force to append '-Hn' argument to show filename and
	// line number of each matched log file entries
	commands = append(commands, "-Hn", config.LocalInfo.LogPath)
	// Do the grep command execution locally
	cmd := exec.Command("grep", commands...)
	// Get the grep command output in stdout pipe
	stdOut, err := cmd.StdoutPipe()
	var stdOutErr []byte
	// A WaitGroup used to sync the line counting command,
	// which are required to execute at last
	var wg sync.WaitGroup
	if flagC == false {
		wg.Add(1)
		// Concurrently execute the line counting grep command
		// Output the matched line total count
		go func() {
			commandsCount := append(commands, "-c")
			cmdCount := exec.Command("grep", commandsCount...)
			stdOutErr, _ = cmdCount.CombinedOutput()
			fmt.Println(string([]byte(stdOutErr)))
			wg.Done()
		}()
	}
	// Print the error if we have, but not exit
	printError(err)
	// Start the command
	err = cmd.Start()
	printError(err)
	for {
		// Each loop we read the pipe in 1024 bytes chunk
		bufferOut := make([]byte, 1024)
		_, errOut := stdOut.Read(bufferOut)

		if errOut != nil {
			if errOut == io.EOF {
				break
			} else {
				printError(errOut)
			}
		}

		//fmt.Printf("%s", string(bufferOut[:lenOut]))
		// Send the grep output to client
		connect.Write(bufferOut)
		// Clean buffer
		bufferOut = nil
	}

	// Wait until the command execution finished
	err = cmd.Wait()
	printError(err)

	if flagC == false {
		// WaitGroup wait the line count cmd finished
		// And send the result to client, client would parse it
		wg.Wait()
		connect.Write([]byte(stdOutErr))
	}

	// Output responses complete
	fmt.Println("[Debug] Output Ended")
	// Close the connection with this client
	connect.Close()
}

// Helper function to print the err in process
func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[ERROR]", err.Error())
	}
}
