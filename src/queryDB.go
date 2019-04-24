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
	_ "time"
)

func queryAuthorizedUsers(host Host, url string) (bool, bool) {
	searchURL := url + "_find/"
	for _, hostname := range host.Hostnames {
		type AddrSelector struct {
			Selector struct {
				Name string `json:"livehost.hostname.name"`
			} `json:"selector"`
		}

		jsonStr := AddrSelector{Selector: struct {
			Name string `json:"livehost.hostname.name"`
		}{Name: hostname.Name}}
		jsonValue, _ := json.Marshal(jsonStr)
		resp, err := http.Post(searchURL, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Error accessing db", err)
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var queryResp FindResponseBody
			err = json.Unmarshal(body, &queryResp)
			if err != nil {
				panic(err)
			}

			if len(queryResp.Docs) > 0 {
				return true, queryResp.Docs[0].Host.Persistent
			}
		}
	}
	return false, false
}

func queryLiveHosts(host Host, url string, checkIn int) bool {
	searchURL := url + "_find/"
	for _, address := range host.Addresses {
		type AddrSelector struct {
			Selector struct {
				Addr string `json:"livehost.ipaddress.addr"`
			} `json:"selector"`
		}

		jsonStr := AddrSelector{Selector: struct {
			Addr string `json:"livehost.ipaddress.addr"`
		}{Addr: address.Addr}}
		jsonValue, _ := json.Marshal(jsonStr)
		resp, err := http.Post(searchURL, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Error accessing databse", err)
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			var queryResp FindResponseBody
			err = json.Unmarshal(body, &queryResp)
			if err != nil {
				panic(err)
			}

			if len(queryResp.Docs) > 0 {
				updateCheckin(queryResp.Docs[0], checkIn, url)
				return true
			}
		}
	}

	return false
}
