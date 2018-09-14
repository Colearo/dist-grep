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
	"strconv"
	"bytes"
)


type Config struct {
	Addresses []string
}

// global variables
var args string
var wg sync.WaitGroup
var mutex sync.Mutex
var total_count int

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
	for index, address := range config.Addresses {
		wg.Add(1)
		go makeRequest(address, index)
	}

	// Wait for all requests to complete.
	wg.Wait()

	// Print total count.
	fmt.Printf("total count: %d.\n", total_count)
	
	// Print processing time.
	end := time.Now()
	fmt.Printf("total time: %.3f seconds.\n", end.Sub(start).Seconds())
}

func makeRequest(address string, index int) {
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


	
	// Read and buffer contents, then print.
	var buf bytes.Buffer
	io.Copy(&buf, conn)
	info := buf.String()
	
	// Retrieve count number from info
	t1 := time.Now()
	info_list := strings.Split(info, ":")
	count_info := info_list[len(info_list) - 1]
	count_info = strings.TrimSpace(count_info)
	count, _ := strconv.Atoi(count_info)
	t2 := time.Now()

	// Synchornize print contents from buffer, and add total_count.
	mutex.Lock()
	total_count += count
	fmt.Print(info)
	fmt.Println("time for retrieve count: ", t2.Sub(t1).Seconds())
	mutex.Unlock()
	


	/*	
	// Read and print concurrently.
	// We ensure the last packet only contains the count information.
	// Variables to collect count info.
	var info, pre_info string
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Print(err)
		}
		pre_info = info
		info = string(buf[:n])
		fmt.Print(pre_info)
	}
	// Move count into to global var count
	//count_str := strings.Split(info, ":")[1]
	//count_str = strings.TrimSpace(count_str)
	//count, _ := strconv.Atoi(count_str)
	// Send count to global variables
	//counts[index] = count
	*/
}
