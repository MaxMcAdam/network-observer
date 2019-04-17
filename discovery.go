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
  "time"
)

func main() {
  checkIn := 0
  url := "http://127.0.0.1:5984/"
  //addrs:=getNetwork()
  //var wg sync.WaitGroup
  for true {
    findChanges(parseNmap(), url, checkIn)
    checkIn++
    time.Sleep(time.Minute)
  }
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

func findChanges(liveHosts []Host, dbURL string, checkIn int) {
  liveHostDBURL := dbURL + "live-hosts/"
  authHostDBURL := dbURL + "auth-hosts/"
  for _,host := range liveHosts {
    exists, docRev, docID := queryLiveHosts(host, liveHostDBURL)
    if exists {
      fmt.Println("host aleady in live host db")
      updateCheckin(docRev, docID, dbURL, checkIn)
    } else {
      authorization, persistence := false, false
      if len(host.Hostnames) > 1{
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

func addHostToLiveHosts(host Host, hostAuthorized bool, hostPersistent bool, url string, checkIn int) {
  var newHostname Hostname
  if len(host.Hostnames) < 1 {
    newHostname = host.Hostnames[0]
  }
  currentTime := time.Now().String()
  newHost := LiveHost{
    IPAddress: host.Addresses[0],
    LiveHostname: newHostname,
    Authorized: hostAuthorized,
    Persistent: hostPersistent,
    LastCheckin: checkIn,
    TimeDiscovered: currentTime,
  }
  jsonStr := map[string]LiveHost{"livehost":newHost}
  jsonValue, _ := json.Marshal(jsonStr)
  fmt.Println(jsonValue)
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
    jsonStr := map[string]Hostname{"selector":searchHostname}
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

func queryLiveHosts(host Host, url string) (bool, string, string){
  searchURL := url + "_find/"
  for _,address := range host.Addresses{
    searchAddress := Address{
      Addr: address.Addr,
      AddrType: address.AddrType,
    }
    jsonStr := map[string]Address{"selector":searchAddress}
    jsonValue, _ := json.Marshal(jsonStr)
    resp, err := http.Post(searchURL, "application/json", bytes.NewBuffer(jsonValue))
    fmt.Println(resp)
    if err!=nil {
      panic(err)
    }
    defer resp.Body.Close()
    var queryResp FindResponseBody
    json.NewDecoder(resp.Body).Decode(&queryResp)
    if queryResp.ExecutionStats.ResultsReturned > 0 {
      return true, queryResp.Docs[0].Rev, queryResp.Docs[0].ID
    }
  }

  return false, "", ""
}

func updateCheckin(docRev string, docID string, url string, checkIn int){
  url = url + docID + "/"
  jsonStr := map[string]interface{}{"lastcheckin":checkIn,"_rev":docRev}
  jsonValue, _ := json.Marshal(jsonStr)
  resp, err := http.NewRequest(http.MethodPut,url, bytes.NewBuffer(jsonValue))
  if err != nil{
    panic(err)
  }
  fmt.Println(resp)
}

type LiveHost struct {
  IPAddress Address `json:"ipaddress"`
  LiveHostname Hostname `json:"livehostname"`
  Authorized bool `json:"authorized"`
  Persistent bool `json:"persistent"`
  LastCheckin int `json:"lastcheckin"`
  TimeDiscovered string `json:"timediscovered"`
  Rev string `json:"_rev"`
  ID string `json:"_id"`
}

type FindResponseBody struct {
  Docs []LiveHost `json:"docs"`
  ExecutionStats ExecStats `json:"execution_stats"`
}

type ExecStats struct {
  TotalKeysExamined int `json:"total_keys_examined"`
  TotalDocsExamined int `json:"total_docs_examined"`
  TotalQuorumDocsExamined int `json:"total_quorum_docs_examined"`
  ResultsReturned int `json:"results_returned"`
  ExecutionTimeMS float32 `json:"execution_time_ms"`
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
