package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	serverAddr := ":5555"
	udpAddr, err := net.ResolveUDPAddr("udp4", serverAddr)
	printError(err)
	connect, err := net.ListenUDP("udp", udpAddr)
	printError(err)
	for {
		handleMsg(connect)
	}
}

func handleMsg(connect *net.UDPConn) {
	var buffer [1024]byte
	_, addr, err := connect.ReadFromUDP(buffer[0:])
	commands := strings.Split(string(buffer[0:]), " ")
	fmt.Println(command)
	if err != nil {
		fmt.Println("Fatal Error")
		return
	}
	cmd := exec.Command("grep", commands)
	stdOutErr, err := cmd.CombinedOutput()
	printError(err)
	fmt.Println(string(stdOutErr[0:]))

	connect.WriteToUDP([]byte(stdOutErr), addr)

}

func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error", err.Error())
		os.Exit(1)
	}
}
