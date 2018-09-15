package rgrep

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"
	"io"
	"io/ioutil"
	"encoding/json"
	"time"
	"sync"
	"strconv"
	"bytes"
	"regexp"
)

// OOP programming. Convenient for unit test.
type Rgrep struct {}

type Config struct {
	Addresses []string
}

// global variables
var args string
var wg sync.WaitGroup
var mutex sync.Mutex
var total_count int
var total_connected_vm int

func (r Rgrep) Launch(input_args string) {
	// Start timer
	start := time.Now()

	// Open local config.json file.
	usr, _ := user.Current()
	usrHome := usr.HomeDir
	configFile, err := os.Open(usrHome + "/go/src/dist-grep/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	// Read json file's contents and cache them to var config.
	configBytes, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(configBytes, &config)
	
	// Store grep arguments if args length is valid.
	// Otherwise use func arguments. For unit testing.
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
	} else {
		args = input_args
	}

	// Send concurrent requests to all servers.
	for index, address := range config.Addresses {
		wg.Add(1)
		go makeRequest(address, index)
	}

	// Wait for all requests to complete.
	wg.Wait()

	// Print total connected VMs.
	fmt.Printf("Total Connected VMs: %d\n", total_connected_vm)

	// Print total count.
	fmt.Printf("Total Counts: %d\n", total_count)
	
	// Print total time.
	end := time.Now()
	fmt.Printf("Total Time: %.3f seconds\n", end.Sub(start).Seconds())
}

func makeRequest(address string, index int) {
	// Notify the WaitGroup after this goroutine complete.
	defer wg.Done()

	// Record whether this machine is connected
	var connected bool

	// Time out needed in order to deal with server failure.
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		fmt.Printf("Failed to connect %s\n", address)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected with %s\n", address)
	connected = true

	// Write arguments to the remote grep server.
	conn.Write([]byte(args))

	// Read and buffer contents.
	var buf bytes.Buffer
	io.Copy(&buf, conn)
	info := buf.String()
	
	// Retrieve count number from info.
	info_list := strings.Split(info, ":")
	count_info := info_list[len(info_list) - 1]
	re := regexp.MustCompile("[0-9]+")
	count_info = re.FindString(count_info)
	count, _ := strconv.Atoi(count_info)

	// Synchornize print contents from buffer, and add total_count.
	mutex.Lock()
	if connected {
		total_connected_vm += 1
	}
	total_count += count
	fmt.Print(info)
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
