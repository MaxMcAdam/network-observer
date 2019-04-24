package main

import (
	"fmt"
	_ "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/exec"
	"strconv"
)

func newAlert(alertType string, alertDevice string) {
	//pubClientOptions := mqtt.NewClientOptions()
	//pubClientOptions.Username = wiotp_org
	//pubClientOptions.Password = wiotp_auth_token
	//pubClient := mqtt.NewClient(pubClientOptions)
	cmdFunction := "mosquitto_pub"
	fmt.Println(cmdFunction)
	alertString := "{" + strconv.Quote("alerttype") + ":" + strconv.Quote(alertType) + "," + strconv.Quote("alertdevice") + ":" + strconv.Quote(alertDevice) + "}"
	cmd := exec.Command(cmdFunction, "-h", os.ExpandEnv("WIOTP_ORG")+".messaging.internetoftings.ibmcloud.com", "-p", "8883", "-i", "d:"+os.ExpandEnv("WIOTP_ORG")+":"+os.ExpandEnv("WIOTP_DEVICE_TYPE")+":"+os.ExpandEnv("WIOTP_DEVICE_ID"), "-u", "use-token-auth", "-P", os.ExpandEnv("WIOTP_DEVICE_TOKEN"), "--capath", "/etc/ssl/certs", "-t", "iot-2/evt/status/fmt/json", "-m", alertString, "-d")

	mosqPubOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error sending alert to wiotp", err)
	}
	fmt.Println(mosqPubOutput)
}
