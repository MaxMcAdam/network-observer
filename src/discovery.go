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
	_ "os"
	"os/exec"
	_ "strconv"
	"strings"
	"sync"
	_ "time"
)

func getNetwork() string {
	network := ""
	netInterfaces, _ := net.Interfaces()
	for _, interf := range netInterfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for _, address := range addrs {
				addrParts := strings.Split(address.String(), "/")
				netAddr := addrParts[0]
				subnet := addrParts[1]
				if netAddr != "127.0.0.1" && netAddr != "::1" && subnet != "64" {
					network = address.String()
				}
			}
		}
	}
	return network
}

func discovery(wg *sync.WaitGroup, ipAddr string, ipv6 bool) []byte {
	cmdFunction := "nmap"
	fmt.Println(cmdFunction, "-sn", "-oX", "-", ipAddr)
	cmd := exec.Command(cmdFunction, "-sn", "-oX", "-", ipAddr)
	if ipv6 {
		cmd = exec.Command(cmdFunction, "-6", "-sn", "-oX", "-", ipAddr)
	}
	outputScan, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", outputScan)

	//scanner := bufio.NewScanner(networkChanges)
	wg.Done()
	return outputScan
}

func parseNmap(scanOutput []byte) []Host {
	var scan NmapRun
	xml.Unmarshal(scanOutput, &scan)
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
				err := addHostToLiveHosts(host, true, persistence, liveHostDBURL, checkIn)
				if err != nil {
					fmt.Println("Error accessing databse")
				} else {
					fmt.Println("Authorized host added")
				}
			} else {
				err := addHostToLiveHosts(host, false, persistence, liveHostDBURL, checkIn)
				if err != nil {
					fmt.Println("Error accessing databse")
				} else {
					fmt.Println("Authorized host added")
				}
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
	var queryResp AllDocsResp
	json.Unmarshal(body, &queryResp)

	for _, row := range queryResp.Rows {
		resp, err := http.Get(baseURL + "live-hosts/" + row.DocId)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var doc Doc
		json.Unmarshal(body, &queryResp)
		if doc.Host.LastCheckin < currentCheckin-4 {
			if doc.Host.Persistent {
				fmt.Println("Persistent host " + doc.Host.LiveHostname.Name + " has dropped")
			} else {
				if doc.Host.LiveHostname.Name != "" {
					fmt.Println("host " + doc.Host.LiveHostname.Name + " has dropped")
				} else {
					fmt.Println("host ", doc.Host.IPAddress.Addr, " has dropped")
				}
			}
			delURL := baseURL + "live-hosts/" + doc.ID + "?rev=" + doc.Rev
			cmd := exec.Command("curl", "-X", "DELETE", delURL, "-H", "Content-Type:application/json")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
