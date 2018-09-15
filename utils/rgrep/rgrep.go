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
type Rgrep struct {
	Args string
	Config Config
	TotalCount int
	TotalConnectedVMs int
	Wg sync.WaitGroup
	Mutex sync.Mutex
	UsrHome string
	// Instance variables for unit test only.
	IsTest bool
	TestLogs []string
}

// A struct to store config information.
type Config struct {
	Addresses []string
}

func (r *Rgrep) Launch(func_args string) {
	// Start timer
	start := time.Now()

	// Open local config.json file.
	usr, _ := user.Current()
	r.UsrHome = usr.HomeDir
	configFile, err := os.Open(r.UsrHome + "/go/src/dist-grep/config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()

	// Read json file's contents and pass them to r.Config.
	configBytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(configBytes, &r.Config)

	// Use func args as command args if os args is empty.
	if len(os.Args) < 2 {
		r.Args = func_args
	} else {
		r.Args = strings.Join(os.Args[1:], " ")
	}

	// Arguments start with "-t" execute test mode.
	if strings.Split(r.Args, " ")[0] == "-t" {
		r.Args = strings.SplitN(r.Args, " ", 2)[1] // Delete "-t" 
		r.IsTest = true
		fmt.Println("Test mode starts...")
		// Make test log info array. For unit test onlt.
		r.TestLogs = make([]string, len(r.Config.Addresses))
	}

	// Send concurrent requests to all servers.
	for index, address := range r.Config.Addresses {
		r.Wg.Add(1)
		go r.MakeRequest(address, index)
	}

	// Wait for all requests to complete.
	r.Wg.Wait()

	// Print total connected VMs.
	fmt.Printf("Total Connected VMs: %d\n", r.TotalConnectedVMs)

	// Print total count.
	fmt.Printf("Total Counts: %d\n", r.TotalCount)
	
	// Print total time.
	end := time.Now()
	fmt.Printf("Total Time: %.3f seconds\n", end.Sub(start).Seconds())

	// Write test log info in unit test mode only.
	if r.IsTest {
		filename := r.UsrHome + "/go/src/dist-grep/test/test_logs/log"
		f, err := os.Create(filename)
		if err != nil {
			fmt.Println("Failed to create test log file.")
			return
		}
		defer f.Close()
		// Write each goroutine's log info in sequence order.
		for i, _ := range r.Config.Addresses {
			f.Write([]byte(r.TestLogs[i]))
		}
	}

	// Reset global count variables after all requests complete.
	r.TotalCount = 0
	r.TotalConnectedVMs = 0
}

func (r *Rgrep) MakeRequest(address string, index int) {
	// Notify the WaitGroup after this goroutine complete.
	defer r.Wg.Done()

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
	conn.Write([]byte(r.Args))

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

	// Synchornize write operation.
	r.Mutex.Lock()

	// Compute total count and total connected VMs.
	if connected {
		r.TotalConnectedVMs += 1
	}
	r.TotalCount += count
	fmt.Print(info)

	// Create test_logs for unit testing.
	if r.IsTest {
		// Generate test log info whose format can be used in unit test.
		info_list := strings.Split(info, "\n")
		var test_log_info string
		for _, entry := range info_list {
			entry_list := strings.Split(entry, ":")
			if len(entry_list) > 2 {
				entry = entry_list[2]
				test_log_info = test_log_info + entry + "\n"
			}
		}
		// Assign this goroutine's test log info to a global string array.
		// The array index equals to it's goroutine ID.
		r.TestLogs[index] = test_log_info
	}

	r.Mutex.Unlock()

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
