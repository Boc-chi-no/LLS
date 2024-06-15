package ip2location

import (
	"github.com/oschwald/geoip2-golang"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"net"
)

type fileData struct {
	CityDB *geoip2.Reader
	IspDB  *geoip2.Reader
}

// IPData Data from the IP library
var IPData fileData

// InitIPData Initialise ip library data into memory
func (f *fileData) InitIPData(geoip2CityBytes []byte, geoip2IspBytes []byte) (rs interface{}) {
	var (
		geoip2CityData []byte
		err            error
	)
	if setting.Cfg.GEOIP2.UseOnlineGEOIP2 {
		log.InfoPrint("Get Online GeoIP2...")
		geoip2CityData, err = GetOnline()
		if err != nil {
			geoip2CityData = geoip2CityBytes
			log.WarnPrint("Get Online GeoIP2 failed: %s", err)
		}
	} else {
		geoip2CityData = geoip2CityBytes
	}
	f.CityDB, err = geoip2.FromBytes(geoip2CityData)
	if err != nil {
		log.PanicPrint("Init GeoIP2 City DB failed: %s", err)
	}
	f.IspDB, err = geoip2.FromBytes(geoip2IspBytes)
	if err != nil {
		log.PanicPrint("Init GeoIP2 City ISP failed: %s", err)
	}

	return true
}

func (f *fileData) Find(ipStr string) (res model.Location) {
	res = model.Location{}

	ip := net.ParseIP(ipStr)

	cityRecord, err := f.CityDB.City(ip)
	if err != nil {
		log.ErrorPrint("GeoIP2 City Parse IP failed: %s", err)
	}

	IspRecord, err := f.IspDB.ISP(ip)
	if err != nil {
		log.ErrorPrint("GeoIP2 ISP Parse IP failed: %s", err)
	}

	res.CountryIsoCode = cityRecord.Country.IsoCode
	res.Country = cityRecord.Country.Names[setting.Cfg.GEOIP2.GEOIP2Language]
	res.City = cityRecord.City.Names[setting.Cfg.GEOIP2.GEOIP2Language]

	res.ISP = IspRecord.ISP
	res.Organization = IspRecord.Organization
	res.AutonomousSystemNumber = IspRecord.AutonomousSystemNumber
	res.AutonomousSystemOrganization = IspRecord.AutonomousSystemOrganization

	return res
}

func Find(ip string) (res model.Location) {
	return IPData.Find(ip)
}
