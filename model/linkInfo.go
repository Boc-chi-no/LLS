package model

import "net/http"

type LinkInfo struct {
	Hash   string      `bson:"hash"`
	IP     string      `bson:"ip"`
	Header http.Header `bson:"header"`
	Location
	UAInfo
	Created int64 `bson:"created"`
}

type UAInfo struct {
	Browser        string `bson:"browser"`
	BrowserVersion string `bson:"browser_version"`
	OS             string `bson:"os"`
	OSVersion      string `bson:"os_version"`
	Device         string `bson:"device"`
}

type Location struct {
	CountryIsoCode               string `bson:"country_iso_code"`
	Country                      string `bson:"country"`
	City                         string `bson:"city"`
	ISP                          string `bson:"isp"`
	Organization                 string `bson:"organization"`
	AutonomousSystemNumber       uint   `bson:"autonomous_system_number"`
	AutonomousSystemOrganization string `bson:"autonomous_system_organization"`
}
