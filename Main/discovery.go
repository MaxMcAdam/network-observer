package main

import (
	_ "bufio"
	_ "encoding/json"
	"encoding/xml"
	"fmt"
	_ "io"
	"io/ioutil"
	"net"
	_ "net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	checkIn := 0
	url := "http://127.0.0.1:5984/"
	//addrs:=getNetwork()
	//var wg sync.WaitGroup
	for true {
		for i := 0; i < 5; i++ {
			findChanges(parseNmap(), url, checkIn)
			checkIn++
			time.Sleep(15 * time.Second)
		}
		findDroppedHosts(url, checkIn)
	}
}

func getNetwork() []string {
	network := make([]string, 0)
	netInterfaces, _ := net.Interfaces()
	for _, interf := range netInterfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for index, address := range addrs {
				if interf.Name == "en0" && index == 1 {
					network = append(network, address.String())
				}
			}
		}
	}
	return network
}

func discovery(wg *sync.WaitGroup, ipAddr string) {
	cmdFunction := "./service.sh"
	cmdArgs := ipAddr
	cmd := exec.Command(cmdFunction, cmdArgs)
	_, err := cmd.StdoutPipe()
	cmd.Start()
	if err != nil {
		fmt.Println("Error creating the standard output pipe")
		panic(err)
	}
	//scanner := bufio.NewScanner(networkChanges)
	wg.Done()
}

func parseNmap() []Host {
	xmlFile, err := os.Open("scan.xml")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var scan NmapRun

	xml.Unmarshal(byteValue, &scan)

	defer xmlFile.Close()
	return scan.Hosts
}

func findChanges(liveHosts []Host, dbURL string, checkIn int) {
	liveHostDBURL := dbURL + "live-hosts/"
	authHostDBURL := dbURL + "auth-hosts/"
	for _, host := range liveHosts {
		exists := queryLiveHosts(host, liveHostDBURL, checkIn)
		if exists {
			fmt.Println("host aleady in live host db")
		} else {
			authorization, persistence := false, false
			if len(host.Hostnames) > 0 {
				fmt.Println("Hostname found")
				authorization, persistence = queryAuthorizedUsers(host, authHostDBURL)
			}
			if authorization {
				addHostToLiveHosts(host, true, persistence, liveHostDBURL, checkIn)
			} else {
				fmt.Println("Unauthorized host added")
				addHostToLiveHosts(host, false, persistence, liveHostDBURL, checkIn)
			}
		}
	}
}

func findDroppedHosts(baseURL string, currentCheckin int) {
	searchURL := baseURL + "live-hosts/_all_docs"
	resp, err := http.Get(searchURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var queryResp FindResponseBody
	json.Unmarshal(body, &queryResp)
	for _, doc := range queryResp.Docs {
		if doc.Host.LastCheckin < currentCheckin-1 {
			fmt.Println("host ", doc.Host.IPAddress.Addr, " has dropped")
		}
	}
}

type LiveHost struct {
	IPAddress      Address  `json:"ipaddress"`
	LiveHostname   Hostname `json:"hostname"`
	Authorized     bool     `json:"authorized"`
	Persistent     bool     `json:"persistent"`
	LastCheckin    int      `json:"lastcheckin"`
	TimeDiscovered string   `json:"timediscovered"`
}

type AuthHost struct {
	AuthHostname Hostname `json:"hostname"`
	Persistent   bool     `json:"persistent"`
	DeviceDesc   string   `json:"devdesc"`
}

type FindResponseBody struct {
	Docs     []Doc  `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
	//ExecutionStats ExecStats `json:"execution_stats"`
}

type Doc struct {
	ID   string   `json:"_id"`
	Rev  string   `json:"_rev"`
	Host LiveHost `json:"livehost"`
}

type ExecStats struct {
	TotalKeysExamined       int     `json:"total_keys_examined"`
	TotalDocsExamined       int     `json:"total_docs_examined"`
	TotalQuorumDocsExamined int     `json:"total_quorum_docs_examined"`
	ResultsReturned         int     `json:"results_returned"`
	ExecutionTimeMS         float32 `json:"execution_time_ms"`
}

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun" json:"nmaprun"`
	Hosts   []Host   `xml:"host" json:"host"`
}

type Host struct {
	XMLName   xml.Name   `xml:"host" json:"host"`
	Status    []Status   `xml:"status" json:"status"`
	Addresses []Address  `xml:"address" json:"address"`
	Hostnames []Hostname `xml:"hostnames>hostname" json:"hostnames"`
}

type Hostname struct {
	XMLName xml.Name `xml:"hostname" json:"hostname"`
	Name    string   `xml:"name,attr" json:"name"`
	Type    string   `xml:"type,attr" json:"type"`
}

type Status struct {
	State     string  `xml:"state,attr" json:"state"`
	Reason    string  `xml:"reason,attr" json:"reason"`
	ReasonTTL float32 `xml:"reason_ttl,attr" json:"reason_ttl"`
}

type Address struct {
	Addr     string `xml:"addr,attr" json:"addr"`
	AddrType string `xml:"addrtype,attr" json:"addrtype"`
}
