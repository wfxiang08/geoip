package geoip

import (
	"fmt"
	. "gen-go/geoip_service"
	"github.com/oschwald/maxminddb-golang"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"io/ioutil"
	"net"
)

var reader *maxminddb.Reader

func IpToGeoData(ipStr string) (*GeoData, error) {
	ip := net.ParseIP(ipStr)

	var city City
	err := reader.Lookup(ip, &city)
	if err != nil {
		return nil, err
	} else {
		geoData := &GeoData{
			CountryName:    city.Country.Names["en"],
			CountryIsoCode: city.Country.IsoCode,
			CityName:       city.City.Names["en"],
			Lat:            fmt.Sprintf("%.7f", city.Location.Latitude),
			Lng:            fmt.Sprintf("%.7f", city.Location.Longitude),
			Timezone:       city.Location.TimeZone,
			Continent: city.Continent.Names["en"],
			ContinentCode: city.Continent.Code,
		}
		if len(city.Subdivisions) > 0 {
			geoData.Province = city.Subdivisions[0].IsoCode
		}
		return geoData, nil
	}
}

// 全内存保存DB数据
func InitMaxMindDb(dbPath string) {
	data, err := ioutil.ReadFile(dbPath)
	if err != nil {
		log.PanicError(err, "Read GeoIP failed")
	}
	reader, err = maxminddb.FromBytes(data)
	if err != nil {
		log.PanicError(err, "Read GeoIP failed")
	}
}
