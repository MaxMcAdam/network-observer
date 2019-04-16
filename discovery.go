package main

import (
	"bufio"
	_ "bytes"
	_ "encoding/json"
	"fmt"
	_ "io"
	_ "net"
	"os/exec"
)

func main() {
	//networkChanges := make(chan int)
	go discovery()

	//msg := <-networkChanges
	//fmt.Println(msg)
}

func discovery() {
	networkMonitor := exec.Command("./test.sh")
	networkChanges, err := networkMonitor.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating the standard output pipe")
		panic(err)
	}
	scanner := bufio.NewScanner(networkChanges)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println("end goroutine")
}

func db_changes(devicesChanged []string) {

}
