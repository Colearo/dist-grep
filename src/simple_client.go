package main 

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	addr_1 = "10.194.16.24:5555"
)

func main() {

	// Store grep arguments.
	args := strings.Join(os.Args[2:], " ")

	conn, err := net.Dial("udp", "10.0.0.197:5555")
	if err != nil {
		// handle error
		panic(err)
	}
	defer conn.Close()

	// write
	conn.Write([]byte(args))
	// read
	buffer := make([]byte, 1024)
	conn.Read(buffer)

	// print the buffer
	fmt.Printf("%s\n", string(buffer[0:]))
	
}