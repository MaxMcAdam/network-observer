package main

import (
	"fmt"
	_ "os"
	"os/exec"
	"strconv"
	"sync"
)

func newAlert(wg *sync.WaitGroup, alertType string, alertDevice string, wiotpenv [4]string) {
	//pubClientOptions := mqtt.NewClientOptions()
	//pubClientOptions.Username = wiotp_org
	//pubClientOptions.Password = wiotp_auth_token
	//pubClient := mqtt.NewClient(pubClientOptions)

	//

	cmdFunction := "mosquitto_pub"
	alertString := "{" + strconv.Quote("alerttype") + ":" + strconv.Quote(alertType) + "," + strconv.Quote("alertdevice") + ":" + strconv.Quote(alertDevice) + "}"
	fmt.Println(cmdFunction, "-h", wiotpenv[0]+".messaging.internetofthings.ibmcloud.com", "-p", "8883", "-i", "d:"+wiotpenv[0]+":"+wiotpenv[1]+":"+wiotpenv[2], "-u", "use-token-auth", "-P", wiotpenv[3], "--capath", "/etc/ssl/certs", "-t", "iot-2/evt/status/fmt/json", "-m", alertString, "-d")
	cmd := exec.Command(cmdFunction, "-h", wiotpenv[0]+".messaging.internetofthings.ibmcloud.com", "-p", "8883", "-i", "d:"+wiotpenv[0]+":"+wiotpenv[1]+":"+wiotpenv[2], "-u", "use-token-auth", "-P", wiotpenv[3], "--capath", "/etc/ssl/certs", "-t", "iot-2/evt/status/fmt/json", "-m", alertString, "-d")

	mosqPubOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error sending alert to wiotp", err)
	}
	fmt.Println(mosqPubOutput)
	wg.Done()
}
