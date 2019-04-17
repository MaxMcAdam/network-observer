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
  url := "http://127.0.0.1:5984/"
  addrs:=getNetwork()
  var wg sync.WaitGroup
  for _,address := range addrs {
    wg.Add(1)
    fmt.Println(address)
    go discovery(&wg, address)
    findChanges(parseNmap(), url)
  }
	wg.Wait()
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
    queryLiveHosts(host, dbURL)
    if(len(host.Hostnames) > 0){
      fmt.Println(host.Hostnames[0].Name)
    } else {
      fmt.Println("no name")
    }
  }
}

func queryLiveHosts(host Host, url string) bool{
  jsonStr := map[string]string{"addr":host.Addresses[0].Addr}
  jsonValue, _ := json.Marshal(jsonStr)
  resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
  if err!=nil{
    return false
  }
  fmt.Println(resp)
  return true
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
