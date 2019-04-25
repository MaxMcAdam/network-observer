package main

import (
	_ "bufio"
	_ "encoding/json"
	"encoding/xml"
	"fmt"
	_ "io"
	"io/ioutil"
	_ "net"
	_ "net/http"
	"os"
	_ "os/exec"
	"strconv"
	_ "strings"
	"sync"
	"time"
)

func main() {
	wiotpenv := [4]string{"lyim2y", "Network-Observer", "Network_observer_vm_test", "net-obs-iotp-token"}
	var wg sync.WaitGroup
	wg.Add(1)
	newAlert(&wg, "test-alert", "test alert device", wiotpenv)
	wg.Wait()
	mock := false
	var err error
	if len(os.Args) > 1 {
		mock, err = strconv.ParseBool(os.Args[1])
		if err != nil {
			mock = false
		}
	}

	checkIn := 0
	url := "http://admin:p4ssw0rd@127.0.0.1:5984/"
	addr := getNetwork()

	missingDB(url)

	for true {
		for i := 0; i < 1; i++ {
			var wg sync.WaitGroup
			var hostList []Host
			if mock {
				xmlFile, err := os.Open("../scan.xml")
				if err != nil {
					fmt.Println(err)
				}
				byteValue, _ := ioutil.ReadAll(xmlFile)
				var scan NmapRun
				xml.Unmarshal(byteValue, &scan)
				defer xmlFile.Close()
				hostList = scan.Hosts
			} else {
				wg.Add(1)
				if addr[0:1] == ":" || addr[1:2] == ":" || addr[2:3] == ":" || addr[3:4] == ":" || addr[4:5] == ":" || addr[5:6] == ":" {
					hostList = append(parseNmap(discovery(&wg, string(addr), true)))
				} else {
					hostList = append(parseNmap(discovery(&wg, string(addr), false)))
				}
				wg.Wait()
			}

			findChanges(hostList, url, checkIn, wiotpenv)
			checkIn++
			time.Sleep(15 * time.Second)
			missingDB(url)
		}
		findDroppedHosts(url, checkIn)
	}
}

func missingDB(url string) {
	for !checkDBConn(url) {
		time.Sleep(5 * time.Second)
		fmt.Println("DB not found")
	}
}
