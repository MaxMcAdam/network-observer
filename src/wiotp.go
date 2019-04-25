package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

func newAlert(wg *sync.WaitGroup, alertType string, alertDevice string) {
	//pubClientOptions := mqtt.NewClientOptions()
	//pubClientOptions.Username = wiotp_org
	//pubClientOptions.Password = wiotp_auth_token
	//pubClient := mqtt.NewClient(pubClientOptions)

	//

	cmdFunction := "mosquitto_pub"
	alertString := "{" + strconv.Quote("alerttype") + ":" + strconv.Quote(alertType) + "," + strconv.Quote("alertdevice") + ":" + strconv.Quote(alertDevice) + "}"
	fmt.Println(cmdFunction, "-h", os.ExpandEnv("WIOTP_ORG")+".messaging.internetoftings.ibmcloud.com", "-p", "8883", "-i", "d:"+os.ExpandEnv("WIOTP_ORG")+":"+os.ExpandEnv("WIOTP_DEVICE_TYPE")+":"+os.ExpandEnv("WIOTP_DEVICE_ID"), "-u", "use-token-auth", "-P", os.ExpandEnv("WIOTP_DEVICE_TOKEN"), "--capath", "/etc/ssl/certs", "-t", "iot-2/evt/status/fmt/json", "-m", alertString, "-d")
	cmd := exec.Command(cmdFunction, "-h", os.ExpandEnv("WIOTP_ORG")+".messaging.internetoftings.ibmcloud.com", "-p", "8883", "-i", "d:"+os.ExpandEnv("WIOTP_ORG")+":"+os.ExpandEnv("WIOTP_DEVICE_TYPE")+":"+os.ExpandEnv("WIOTP_DEVICE_ID"), "-u", "use-token-auth", "-P", os.ExpandEnv("WIOTP_DEVICE_TOKEN"), "--capath", "/etc/ssl/certs", "-t", "iot-2/evt/status/fmt/json", "-m", alertString, "-d")

	mosqPubOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error sending alert to wiotp", err)
	}
	fmt.Println(mosqPubOutput)
	wg.Done()
}
