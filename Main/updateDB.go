package main

import (
	_ "bufio"
	"bytes"
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	_ "io"
	_ "io/ioutil"
	_ "net"
	"net/http"
	_ "os"
	"os/exec"
	"strconv"
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

  ipAddressStr := "{" + strconv.Quote("addr") + ":" + strconv.Quote(docToRev.Host.IPAddress.Addr) + "," + strconv.Quote("addrtype") + ":" + strconv.Quote(docToRev.Host.IPAddress.AddrType)"}"
  hostnameStr := "{" + strconv.Quote("name") + ":" + strconv.Quote(docToRev.Host.LiveHostname.Name) + "," + strconv.Quote("type") + ":" + strconv.Quote(docToRev.Host.LiveHostname.Type)"}"
  liveHostStr := "{" + strconv.Quote("ipaddress") + ":" + ipAddressStr + "," + strconv.Quote("hostname") + ":" + hostnameStr + "," + strconv.Quote("authorized") + ":" + strconv.Quote(docToRev.Host.Authorized) + "," + strconv.Quote("persistent") + ":" + strconv.Quote(docToRev.Host.Persistent) + "," + strconv.Quote("lastcheckin") + ":" + strconv.Quote(docToRev.Host.LastCheckin) + "," + strconv.Quote("timediscovered") + ":" + strconv.Quote(docToRev.Host.TimeDiscovered) + "}"
	strDocToRev := "'{" + strconv.Quote("id") + ":" + strconv.Quote(docToRev.ID) + "," + strconv.Quote("rev") + ":" + strconv.Quote(docToRev.Rev) + "," + strconv.Quote("livehost") + ":" + liveHostStr + "}'")
	cmd := exec.Command("curl", "-X", "PUT", updateURL, "-H", "Content-Type:application/json", "-d", strDocToRev)
	output, _ := cmd.CombinedOutput()
	//var strOutput
	fmt.Println(string(output))

	//fmt.Println(printable)
	//resp, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(jsonValue))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//defer resp.Body.Close()
	//fmt.Println("http PUT Response Header: ", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//var queryResp FindResponseBody
	//fmt.Println("http PUT Response Body: ", queryResp)
	//json.Unmarshal(body, &queryResp)
}
