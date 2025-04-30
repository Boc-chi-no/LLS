package uap

import (
	"linkshortener/log"
	"linkshortener/model"
	"net/http"

	"github.com/ua-parser/uap-go/uaparser"
)

var parser *uaparser.Parser

func InitUap(data []byte) {
	var err error
	parser, err = uaparser.NewFromBytes(data)
	if err != nil {
		log.PanicPrint("Loading UAInfo failed: %s", err)
	}
}

func Parse(req http.Header) model.UAInfo {
	uap := model.UAInfo{}
	if value, isExist := req["User-Agent"]; isExist {
		client := parser.Parse(value[0])
		uap.Device = client.Device.ToString()
		uap.OS = client.Os.Family
		uap.OSVersion = client.Os.ToVersionString()
		uap.Browser = client.UserAgent.Family
		uap.BrowserVersion = client.UserAgent.ToVersionString()
	}
	return uap
}
