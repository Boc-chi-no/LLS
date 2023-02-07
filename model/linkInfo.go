package model

import "net/http"

type LinkInfo struct {
	Hash              string         `bson:"hash"`
	IP                string         `bson:"ip"`
	Header            http.Header    `bson:"header"`
	Location
	UAInfo
	Created           int64          `bson:"created"`

}

type UAInfo struct {
	Browser           string         `bson:"browser"`
	BrowserVersion    string         `bson:"browser_version"`
	OS                string         `bson:"os"`
	OSVersion         string         `bson:"os_version"`
	Device            string         `bson:"device"`
}

type Location struct {
	Country           string         `bson:"country"`
	Area              string         `bson:"area"`
}
