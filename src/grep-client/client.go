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
	"sync"
)


type Config struct {
	Addresses []string
}

// global variables
var args string
var wg sync.WaitGroup

func main() {
	// Start timer
	start := time.Now()

	// Open local config.json file.
	configFile, err := os.Open("../../config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	// Read json file's contents and cache them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(configBytes, &config)
	
	// Store grep arguments.
	args = strings.Join(os.Args[1:], " ")

	// Send concurrent requests to all servers.
	for _, address := range config.Addresses {
		wg.Add(1)
		go makeRequest(address)
	}

	// Wait for all requests to complete.
	wg.Wait()
	end := time.Now()
	fmt.Printf("total time: %.3f seconds.\n", end.Sub(start).Seconds())
}

func makeRequest(address string) {
	// Notify the WaitGroup after this goroutine complete.
	defer wg.Done()

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

	// Read and print concurrently.
	// We ensure the last packet only contains the count information.
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		//fmt.Printf("\nNNNNNNNNNNNNNNNNNNNNN: %d\n", n)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Print(err)
		}
		fmt.Print(string(buf[:n]))
	}
}
