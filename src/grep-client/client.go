package main 

import (
	"fmt"
	"net"
	"os"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
	"time"
)


type Config struct {
	Addresses []string
}

// global var
var args string

func main() {
	// Open local config.json file.
	configFile, err := os.Open("../../config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	// Read json file's contents and pass them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(configBytes, &config)
	
	// Store grep arguments.
	args = strings.Join(os.Args[1:], " ")

	for i := 0; i < len(config.Addresses); i++ {
		go makeRequest(config.Addresses[i])
	}

	select{}
}

func makeRequest(address string) {

	// Time out needed in order to deal with server failure.
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		fmt.Printf("Failed to connect %s\n", address)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected with %s\n", address)

	// Write arguments to the remote grep server.
	conn.Write([]byte(args))

	//Read and print concurrently.
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				// Handle input contents.
				fmt.Println(string(buf[:n]))
				break
			}
			fmt.Println(err)
		}
		// Handle input contents.
		fmt.Print(string(buf[:n]))
	}
}
