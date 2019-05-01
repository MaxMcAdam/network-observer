package main

import (
	_ "bufio"
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	_ "io"
	"io/ioutil"
	_ "net"
	_ "net/http"
	_ "os"
	_ "os/exec"
	"strconv"
	_ "strings"
	"sync"
	"time"
)

func main() {
	userVars := getUserVars()
	wiotpenv := [4]string{userVars.WiotpOrg, userVars.WiotpDeviceType, userVars.WiotpDeviceID, userVars.WiotpDeviceToken}
	urlRam := "https://" + userVars.DbAdminUser + userVars.DbAdminPWstring + "@" + userVars.DbURL + ":" + "5984"
	urlDisk := "https://" + userVars.DbAdminUser + userVars.DbAdminPWstring + "@" + userVars.DbURL + ":" + "5985"
	pauseLength, _ := strconv.ParseInt(userVars.PauseBetweenNmapS, 0, 64)
	pausesBeforeSync, _ := strconv.ParseInt(userVars.PausesBeforeSync, 0, 64)

	checkIn := 0

	addr := getNetwork()

	if checkDBConn(urlDisk) {
		syncDB(0)
	} else {
		initDB()
	}

	for true {
		for i := 0; i < int(pausesBeforeSync); i++ {
			var wg sync.WaitGroup
			var hostList []Host

			wg.Add(1)
			if addr[0:1] == ":" || addr[1:2] == ":" || addr[2:3] == ":" || addr[3:4] == ":" || addr[4:5] == ":" || addr[5:6] == ":" {
				hostList = append(parseNmap(discovery(&wg, string(addr), true)))
			} else {
				hostList = append(parseNmap(discovery(&wg, string(addr), false)))
			}
			wg.Wait()

			findChanges(hostList, urlRam, checkIn, wiotpenv)
			checkIn++
			time.Sleep(time.Duration(pauseLength) * time.Second)

			missingDB(urlRam)
		}
		syncDB(1)
	}
}

func missingDB(url string) {
	for !checkDBConn(url) {
		time.Sleep(5 * time.Second)
		fmt.Println("DB not found")
	}
}

func getUserVars() InitVars {
	userVarsJSON, err := ioutil.ReadFile("/home/edgenode/Documents/network-observer/src/service-vars.json")
	if err != nil {
		fmt.Println("Error reading user variables", err)
	}
	var userVars InitVars
	json.Unmarshal(userVarsJSON, &userVars)
	return userVars
}
