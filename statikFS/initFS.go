package statikFS

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/rakyll/statik/fs"
	"linkshortener/i18n"
	"linkshortener/lib/ip2location"
	"linkshortener/lib/uap"
	"linkshortener/log"
	"net/http"
)

var StatikFS http.FileSystem
var CaptchaFont *truetype.Font

func InitFs() {
	var err error
	StatikFS, err = fs.New()

	if err != nil {
		log.PanicPrint("Init StatikFS failed", err)
	}
}

func InitFont() {
	fontBytes, err := fs.ReadFile(StatikFS, "/statik/arphic.ttf")

	if err != nil {
		log.PanicPrint("Init Font failed", err)
	}

	CaptchaFont, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.PanicPrint("Init Font failed", err)
	}
}

func InitUap() {
	uapBytes, err := fs.ReadFile(StatikFS, "/statik/uaparser.yaml")
	if err != nil {
		log.PanicPrint("Init UAInfo failed", err)
	}

	uap.InitUap(uapBytes)
}

func InitIPData() {
	ipDataBytes, err := fs.ReadFile(StatikFS, "/statik/qqwry.dat")
	if err != nil {
		log.PanicPrint("Init IPData failed", err)
	}

	ip2location.IPData.InitIPData(ipDataBytes)
}

func InitI18n() {
	jpBytes, err := fs.ReadFile(StatikFS, "/statik/lang/ja-JP.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(ja-JP) exception: %s", err)
	}
	cnBytes, err := fs.ReadFile(StatikFS, "/statik/lang/zh-CN.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(zh-CN) exception: %s", err)
	}
	usBytes, err := fs.ReadFile(StatikFS, "/statik/lang/en-US.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(en-US) exception: %s", err)
	}
	i18n.InitI18n(jpBytes, cnBytes, usBytes)
}
