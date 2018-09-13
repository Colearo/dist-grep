package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	serverAddr := ":5555"
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
	commands = append(commands, "-Hn", "/Users/colearolu/Downloads/logs/vm1.log")
	cmd := exec.Command("grep", commands...)
	stdOut, err := cmd.StdoutPipe()
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

	commands = append(commands, "-c")
	stdOutErr, _ := cmd.CombinedOutput()
	connect.Write([]byte(stdOutErr))

	fmt.Println("[Debug] Output Ended")
	connect.Close()
}

func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]", err.Error())
		os.Exit(1)
	}
}
