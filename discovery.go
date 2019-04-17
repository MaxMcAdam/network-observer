package main

import (
	"bufio"
	"bytes"
	_ "encoding/json"
	"fmt"
	_ "io"
	_ "net"
	"os/exec"
  "os"
	"sync"
  "net"
  "encoding/xml"
  "io/ioutil"
  "net/http"
  "encoding/json"
)

func main() {
  url := "http://127.0.0.1:5984/live-hosts/"
  //addrs:=getNetwork()
  //var wg sync.WaitGroup
  //for _,address := range addrs {
    //wg.Add(1)
    //fmt.Println(address)
    //go discovery(&wg, address)
    findChanges(parseNmap(), url)
  //}
	//wg.Wait()
}

func getNetwork() []string{
  network := make([]string, 0)
  netInterfaces, _ := net.Interfaces()
  for _, interf := range netInterfaces{
    if addrs, err := interf.Addrs(); err == nil{
      for index,address := range addrs {
        if(interf.Name == "en0" && index == 1){
          network = append(network, address.String())
        }
      }
    }
  }
  return network
}

func discovery(wg *sync.WaitGroup, ipAddr string) {
	fmt.Println("start goroutine")
  cmdFunction := "./service.sh"
  cmdArgs := ipAddr
	cmd := exec.Command(cmdFunction, cmdArgs)
	networkChanges, err := cmd.StdoutPipe()
	cmd.Start()
	if err != nil {
		fmt.Println("Error creating the standard output pipe")
		panic(err)
	}
	scanner := bufio.NewScanner(networkChanges)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	fmt.Println("end goroutine")
	wg.Done()
}

func parseNmap() []Host{
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

func findChanges(liveHosts []Host, dbURL string) {
  for _,host := range liveHosts {
    exists := queryLiveHosts(host, dbURL)
    if exists {
      fmt.Println("host aleady in live host db")
    } else {
      if len(host.Hostnames) > 1{
        authorization, persistence := queryAuthorizedUsers(host, dbURL)
        if authorization {
          addHostToLiveHosts(host, true, persistence, dbURL)
        } else {
          fmt.Println("Unauthorized host added")
          addHostToLiveHosts(host, false, persistence, dbURL)
        }
      }
      fmt.Println(host.Addresses[0].Addr, " not found")
    }
  }
}

func addHostToLiveHosts(host Host, hostAuthorized bool, hostPersistent bool, url string) {
  var newHostname Hostname
  if len(host.Hostnames) < 1 {
    newHostname = host.Hostnames[0]
  }
  newHost := LiveHost{
    IPAddress: host.Addresses[0],
    LiveHostname: newHostname,
    Authorized: hostAuthorized,
    Persistent: hostPersistent,
    LastCheckin: 0,
  }
  jsonStr := map[string]LiveHost{"livehost":newHost}
  jsonValue, _ := json.Marshal(jsonStr)
  resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
  if err != nil {
    panic(err)
  }
  fmt.Println(resp)
}

func queryAuthorizedUsers(host Host, url string) (bool, bool){
  searchURL := url + "_find/"
  for _,hostname := range host.Hostnames{
    searchHostname := Hostname{
      Name: hostname.Name,
      Type: hostname.Type,
    }
    jsonStr := map[string]Hostname{"hostname":searchHostname}
    jsonValue, _ := json.Marshal(jsonStr)
    resp, err := http.Post(searchURL, "application/json", bytes.NewBuffer(jsonValue))
    if err != nil{
      panic(err)
    }
    if resp.StatusCode != 404 {
      return true, false
    }
  }
  return false, false
}

func queryLiveHosts(host Host, url string) bool{
  searchURL := url + "_find/"
  for _,address := range host.Addresses{
    searchAddress := Address{
      Addr: address.Addr,
      AddrType: address.AddrType,
    }
    jsonStr := map[string]Address{"ip":searchAddress}
    jsonValue, _ := json.Marshal(jsonStr)
    resp, err := http.Post(searchURL, "application/json", bytes.NewBuffer(jsonValue))
    if err!=nil {
      panic(err)
    }
    if resp.StatusCode != 404 {
      return true
    }
  }

  return false
}

type LiveHost struct {
  IPAddress Address `json:"ipaddress"`
  LiveHostname Hostname `json:"livehostname"`
  Authorized bool `json:"authorized"`
  Persistent bool `json:"persistent"`
  LastCheckin int `json:"lastcheckin"`
}

type NmapRun struct {
  XMLName xml.Name `xml:"nmaprun" json:"nmaprun"`
  Hosts []Host `xml:"host" json:"host"`
}

type Host struct {
  XMLName xml.Name `xml:"host" json:"host"`
  Status []Status `xml:"status" json:"status"`
  Addresses []Address `xml:"address" json:"address"`
  Hostnames []Hostname `xml:"hostnames>hostname" json:"hostnames"`
}

type Hostname struct {
  Name string `xml:"name,attr" json:"name"`
  Type string `xml:"type,attr" json:"type"`
}

type Status struct {
  State string `xml:"state,attr" json:"state"`
  Reason string `xml:"reason,attr" json:"reason"`
  ReasonTTL float32 `xml:"reason_ttl,attr" json:"reason_ttl"`
}

type Address struct {
  Addr string `xml:"addr,attr" json:"addr"`
  AddrType string `xml:"addrtype,attr" json:"addrtype"`
}
