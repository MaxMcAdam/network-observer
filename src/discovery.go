package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"sync"
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
				if netAddr != "127.0.0.1" && netAddr != "::1" && subnet != "64" && subnet != "16" {
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

func findChanges(liveHosts []Host, dbURL string, checkIn int, wiotpenv [4]string) {
	liveHostDBURL := dbURL + "live-hosts/"
	authHostDBURL := dbURL + "auth-hosts/"
	for _, host := range liveHosts {
		exists := queryLiveHosts(host, liveHostDBURL, checkIn)
		if !exists {
			var wg sync.WaitGroup
			authorization, persistence := false, false
			hostNamed := len(host.Hostnames) > 0
			if hostNamed {
				authorization, persistence = queryAuthorizedUsers(host, authHostDBURL)
			}
			if authorization {
				err := addHostToLiveHosts(host, true, persistence, liveHostDBURL, checkIn)
				if err != nil {
					fmt.Println("Error accessing databse")
				} else {
					fmt.Println("Authorized host added")
					wg.Add(1)
					newAlert(&wg, "new-auth-host", host.Hostnames[0].Name, wiotpenv)
				}
			} else {
				err := addHostToLiveHosts(host, false, persistence, liveHostDBURL, checkIn)
				if err != nil {
					fmt.Println("Error accessing databse")
				} else {
					fmt.Println("Unauthorized host added")
					if hostNamed {
						wg.Add(1)
						newAlert(&wg, "new-unauth-host", host.Hostnames[0].Name, wiotpenv)
					} else {
						wg.Add(1)
						newAlert(&wg, "new-unauth-host", host.Addresses[0].Addr, wiotpenv)
					}
				}
			}
			wg.Wait()
			findDroppedHosts(dbURL, checkIn, wiotpenv)
		}
	}
}

func findDroppedHosts(baseURL string, currentCheckin int, wiotpenv [4]string) {
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
		json.Unmarshal(body, &doc)
		if doc.Host.LastCheckin < currentCheckin-4 {
			var wg sync.WaitGroup
			if doc.Host.Persistent {
				fmt.Println("Persistent host " + doc.Host.LiveHostname.Name + " has dropped")
				wg.Add(1)
				newAlert(&wg, "persistent-host-dropped", doc.Host.LiveHostname.Name, wiotpenv)
			} else {
				if doc.Host.LiveHostname.Name != "" {
					fmt.Println("host " + doc.Host.LiveHostname.Name + " has dropped")
					wg.Add(1)
					newAlert(&wg, "host-dropped", doc.Host.LiveHostname.Name, wiotpenv)
				} else {
					fmt.Println("host ", doc.Host.IPAddress.Addr, " has dropped")
					wg.Add(1)
					newAlert(&wg, "host-dropped", doc.Host.IPAddress.Addr, wiotpenv)
				}
			}
			wg.Wait()
			delURL := baseURL + "live-hosts/" + doc.ID + "?rev=" + doc.Rev
			cmd := exec.Command("curl", "-X", "DELETE", delURL, "-H", "Content-Type:application/json")
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
