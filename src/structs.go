package main

import (
	_ "encoding/json"
	"encoding/xml"
)

type LiveHost struct {
	IPAddress      Address  `json:"ipaddress"`
	LiveHostname   Hostname `json:"hostname"`
	Authorized     bool     `json:"authorized"`
	Persistent     bool     `json:"persistent"`
	LastCheckin    int      `json:"lastcheckin"`
	TimeDiscovered string   `json:"timediscovered"`
}

type AuthHost struct {
	AuthHostname Hostname `json:"hostname"`
	Persistent   bool     `json:"persistent"`
	DeviceDesc   string   `json:"devdesc"`
}

type FindResponseBody struct {
	Docs     []Doc  `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
	//ExecutionStats ExecStats `json:"execution_stats"`
}

type Doc struct {
	ID   string   `json:"_id"`
	Rev  string   `json:"_rev"`
	Host LiveHost `json:"livehost"`
}

type ExecStats struct {
	TotalKeysExamined       int     `json:"total_keys_examined"`
	TotalDocsExamined       int     `json:"total_docs_examined"`
	TotalQuorumDocsExamined int     `json:"total_quorum_docs_examined"`
	ResultsReturned         int     `json:"results_returned"`
	ExecutionTimeMS         float32 `json:"execution_time_ms"`
}

type NmapRun struct {
	XMLName xml.Name `xml:"nmaprun" json:"nmaprun"`
	Hosts   []Host   `xml:"host" json:"host"`
}

type Host struct {
	XMLName   xml.Name   `xml:"host" json:"host"`
	Status    []Status   `xml:"status" json:"status"`
	Addresses []Address  `xml:"address" json:"address"`
	Hostnames []Hostname `xml:"hostnames>hostname" json:"hostnames"`
}

type Hostname struct {
	XMLName xml.Name `xml:"hostname" json:"hostname"`
	Name    string   `xml:"name,attr" json:"name"`
	Type    string   `xml:"type,attr" json:"type"`
}

type Status struct {
	State     string  `xml:"state,attr" json:"state"`
	Reason    string  `xml:"reason,attr" json:"reason"`
	ReasonTTL float32 `xml:"reason_ttl,attr" json:"reason_ttl"`
}

type Address struct {
	Addr     string `xml:"addr,attr" json:"addr"`
	AddrType string `xml:"addrtype,attr" json:"addrtype"`
}

type AllDocsResp struct {
	TotalRow int          `json:"total_rows"`
	Offset   int          `json:"offset"`
	Rows     []AllDocsRow `json:"rows"`
}

type AllDocsRow struct {
	DocId  string `json:"id"`
	DocKey string `json:"key"`
	DocRev string `json:"value.rev"`
}

type InitVars struct {
	WiotpDeviceToken  string `json:"WIOTP_DEVICE_TOKEN"`
	WiotpDeviceID     string `json:"WIOTP_DEVICE_ID"`
	WiotpDeviceType   string `json:"WIOTP_DEVICE_TYPE"`
	WiotpOrg          string `json:"WIOTP_ORG"`
	DbURL             string `json:"DB_URL"`
	DbAdminUser       string `json:"DB_ADMIN_USERNAME"`
	DbAdminPWstring   string `json:"DB_ADMIN_PW"`
	PauseBetweenNmapS string `json:"PAUSE_BETWEEN_NMAP_S"`
	PausesBeforeSync  string `json:"PAUSES_BEFORE_SYNC"`
}
