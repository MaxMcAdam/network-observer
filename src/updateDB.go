package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func addHostToLiveHosts(host Host, hostAuthorized bool, hostPersistent bool, url string, checkIn int) error {
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
		return err
	}
	return nil
}

func updateCheckin(docToRev Doc, checkIn int, url string) {
	updateURL := url + docToRev.ID + "/"
	docToRev.Host.LastCheckin = checkIn

	ipAddressStr := "{" + strconv.Quote("addr") + ":" + strconv.Quote(docToRev.Host.IPAddress.Addr) + "," + strconv.Quote("addrtype") + ":" + strconv.Quote(docToRev.Host.IPAddress.AddrType) + "}"
	hostnameStr := "{" + strconv.Quote("name") + ":" + strconv.Quote(docToRev.Host.LiveHostname.Name) + "," + strconv.Quote("type") + ":" + strconv.Quote(docToRev.Host.LiveHostname.Type) + "}"
	liveHostStr := "{" + strconv.Quote("ipaddress") + ":" + ipAddressStr + "," + strconv.Quote("hostname") + ":" + hostnameStr + "," + strconv.Quote("authorized") + ":" + fmt.Sprintf("%v", docToRev.Host.Authorized) + "," + strconv.Quote("persistent") + ":" + fmt.Sprintf("%v", docToRev.Host.Persistent) + "," + strconv.Quote("lastcheckin") + ":" + fmt.Sprintf("%v", docToRev.Host.LastCheckin) + "," + strconv.Quote("timediscovered") + ":" + strconv.Quote(docToRev.Host.TimeDiscovered) + "}"
	strDocToRev := "{" + strconv.Quote("_id") + ":" + strconv.Quote(docToRev.ID) + "," + strconv.Quote("_rev") + ":" + strconv.Quote(docToRev.Rev) + "," + strconv.Quote("livehost") + ":" + liveHostStr + "}"

	cmd := exec.Command("curl", "-X", "PUT", updateURL, "-H", "Content-Type:application/json", "-d", strDocToRev)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func syncDB(sourceURL string, targetURL string) {
	jsonStr := map[string]string{"source": sourceURL, "target": targetURL}
	jsonValue, _ := json.Marshal(jsonStr)
	_, err := http.Post(targetURL+"_replicate", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Error syncing DB's: ", err)
	}
}

func initDB(url string) {
	reqDBs := []string{"_users", "_replicator", "_global_changes", "auth-hosts", "live-hosts"}
	for _, db := range reqDBs {
		fmt.Println("curl", "-X", "PUT", url+db, "-H", "Content-Type:application/json")
		cmd := exec.Command("curl", "-X", "PUT", url+db, "-H", "Content-Type:application/json")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error creating the databases", err)
		}
	}
}
