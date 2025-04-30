package fs

import (
	"linkshortener/i18n"
	"linkshortener/lib/ip2location"
	"linkshortener/lib/uap"
	"linkshortener/log"
	"net/http"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/rakyll/statik/fs"
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
	fontBytes, err := fs.ReadFile(StatikFS, "/resources/arphic.ttf")

	if err != nil {
		log.PanicPrint("Init Font failed", err)
	}

	CaptchaFont, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.PanicPrint("Init Font failed", err)
	}
}

func InitUap() {
	uapBytes, err := fs.ReadFile(StatikFS, "/resources/uaparser.yaml")
	if err != nil {
		log.PanicPrint("Init UAInfo failed", err)
	}

	uap.InitUap(uapBytes)
}

func InitIPData() {
	geoip2CityBytes, err := fs.ReadFile(StatikFS, "/resources/GeoIP2-City.mmdb")
	if err != nil {
		log.PanicPrint("Init IPData-City failed", err)
	}
	geoip2IspBytes, err := fs.ReadFile(StatikFS, "/resources/GeoIP2-ISP.mmdb")
	if err != nil {
		log.PanicPrint("Init IPData-ISP failed", err)
	}

	ip2location.IPData.InitIPData(geoip2CityBytes, geoip2IspBytes)
}

func InitI18n() {
	jpBytes, err := fs.ReadFile(StatikFS, "/resources/lang/ja-JP.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(ja-JP) exception: %s", err)
	}
	cnBytes, err := fs.ReadFile(StatikFS, "/resources/lang/zh-CN.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(zh-CN) exception: %s", err)
	}
	usBytes, err := fs.ReadFile(StatikFS, "/resources/lang/en-US.json")
	if err != nil {
		log.PanicPrint("Loading embedded language pack(en-US) exception: %s", err)
	}
	i18n.InitI18n(jpBytes, cnBytes, usBytes)
}
