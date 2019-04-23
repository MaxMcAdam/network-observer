package main

import (
	_ "bufio"
	"encoding/json"
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
	_ "time"
)

func main() {
	//checkIn := 0
	//url := "http://127.0.0.1:5984/"
	addrs := getNetwork()
	var wg sync.WaitGroup
	discovery(&wg, addrs[0], "scan.xml")
	//for true {
	//	for i := 0; i < 5; i++ {
	//		findChanges(parseNmap(), url, checkIn)
	//		checkIn++
	//		time.Sleep(15 * time.Second)
	//	}
	//	findDroppedHosts(url, checkIn)
	//}
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

func discovery(wg *sync.WaitGroup, ipAddr string, scanFile string) {
	cmdFunction := "nmap"
	cmd := exec.Command(cmdFunction, "-sN", "-oX", "-", ipAddr, ">", scanFile)
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
		if !exists {
			authorization, persistence := false, false
			if len(host.Hostnames) > 0 {
				authorization, persistence = queryAuthorizedUsers(host, authHostDBURL)
			}
			if authorization {
				addHostToLiveHosts(host, true, persistence, liveHostDBURL, checkIn)
				fmt.Println("Authorized host added")
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
		if doc.Host.LastCheckin < currentCheckin-2 {
			if doc.Host.Persistent {
				fmt.Println("Persistent host " + doc.Host.LiveHostname.Name + " has dropped")
			} else {
				if doc.Host.LiveHostname.Name != "" {
					fmt.Println("host " + doc.Host.LiveHostname.Name + " has dropped")
				} else {
					fmt.Println("host ", doc.Host.IPAddress.Addr, " has dropped")
				}
			}
		}
	}
}
