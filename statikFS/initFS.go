package statikFS

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/rakyll/statik/fs"
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
	ipDataBytes, err := fs.ReadFile(StatikFS, "/resources/qqwry.dat")
	if err != nil {
		log.PanicPrint("Init IPData failed", err)
	}

	ip2location.IPData.InitIPData(ipDataBytes)
}
