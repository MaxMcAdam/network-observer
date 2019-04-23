package main

import (
	_ "bufio"
	"bytes"
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	_ "io"
	"io/ioutil"
	_ "net"
	"net/http"
	_ "os"
	_ "os/exec"
	_ "sync"
	"time"
)

func addHostToLiveHosts(host Host, hostAuthorized bool, hostPersistent bool, url string, checkIn int) {
	var newHostname Hostname
	if len(host.Hostnames) > 0 {
		newHostname = host.Hostnames[0]
	}
	currentTime := time.Now().String()
	newHost := LiveHost{
		IPAddress:      host.Addresses[0],
		LiveHostname:   newHostname,
		Authorized:     hostAuthorized,
		Persistent:     hostPersistent,
		LastCheckin:    checkIn,
		TimeDiscovered: currentTime,
	}
	jsonStr := map[string]LiveHost{"livehost": newHost}
	jsonValue, _ := json.Marshal(jsonStr)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
}

func updateCheckin(docToRev Doc, checkIn int, url string) {
	updateURL := url + docToRev.ID + "/"
	docToRev.Host.LastCheckin = checkIn
	jsonValue, _ := json.Marshal(docToRev)
	//var printable Doc
	//_ = json.Unmarshal(jsonValue, &printable)
	//fmt.Println("\n", printable)

	//fmt.Println()
	//strDocToRev := "''" + string(jsonValue) + "'"
	//cmd := exec.Command("curl", "-X", "PUT", updateURL, "-H", "Content-Type:application/json", "-d", strDocToRev)
	//output, _ := cmd.CombinedOutput()
	//var strOutput
	//fmt.Println(string(output))
	//cmd.Run()

	//fmt.Println(printable)
	resp, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	fmt.Println("http PUT Response Header: ", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	var queryResp FindResponseBody
	fmt.Println("http PUT Response Body: ", queryResp)
	json.Unmarshal(body, &queryResp)
}