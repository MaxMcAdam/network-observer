package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
	"time"
)

func main() {
	userVars := getUserVars()
	wiotpenv := [4]string{userVars.WiotpOrg, userVars.WiotpDeviceType, userVars.WiotpDeviceID, userVars.WiotpDeviceToken}
	urlRam := "http://" + userVars.DbAdminUser + ":" + userVars.DbAdminPWstring + "@" + userVars.DbURL + ":" + "5984" + "/"
	urlDisk := "http://" + userVars.DbAdminUser + ":" + userVars.DbAdminPWstring + "@" + userVars.DbURL + ":" + "5985" + "/"
	pauseLength, _ := strconv.ParseInt(userVars.PauseBetweenNmapS, 0, 64)
	pausesBeforeSync, _ := strconv.ParseInt(userVars.PausesBeforeSync, 0, 64)
	addr := userVars.LanNetwork
	if addr == "" {
		addr = getNetwork()
		if addr == "" {
			panic("Lan Network not found. Please provide network.")
		}
	}
	checkIn := 0

	if checkDBConn(urlDisk) {
		fmt.Println("Persistent db's found")
		syncDB(urlDisk, urlRam)
	} else {
		fmt.Println("Initializing ram db's")
		initDB(urlRam)
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
			if len(hostList) == 0 {
				fmt.Println("No live hosts found on network")
			}

			findChanges(hostList, urlRam, checkIn, wiotpenv)
			checkIn++
			time.Sleep(time.Duration(pauseLength) * time.Second)

			missingDB(urlRam)
		}
		syncDB(urlRam, urlDisk)
	}
}

func missingDB(url string) {
	for !checkDBConn(url) {
		time.Sleep(5 * time.Second)
		fmt.Println("DB not found")
	}
}

func getUserVars() InitVars {
	userVarsJSON, err := ioutil.ReadFile("service-vars.json")
	if err != nil {
		fmt.Println("Error reading user variables", err)
	}
	var userVars InitVars
	json.Unmarshal(userVarsJSON, &userVars)
	return userVars
}
