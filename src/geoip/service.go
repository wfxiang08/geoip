package geoip

import (
	"fmt"
	. "geoip/services"
	"github.com/oschwald/maxminddb-golang"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"github.com/wfxiang08/thrift_rpc_base/rpcthrift/services"
	"io/ioutil"
	"net"
	"time"
	"encoding/json"
)

// 只是为了方便定位
const (
	kErrorCodeNotFound = 1111
)

var reader *maxminddb.Reader

func microseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

// 内部函数不要随便返回一般的Error, 统一封装成为: RpcException
// 注意陷阱: https://golang.org/doc/faq#nil_error
//
func IpToGeoData(ipStr string) (*GeoData, error) {
	start := microseconds()
	ip := net.ParseIP(ipStr)
	var city City
	err := reader.Lookup(ip, &city)

	if err != nil {
		log.Debugf("WARNGIN: IP: %s, no result found, elapsed: %.3fms", ipStr, float64(microseconds() - start) * 0.001)
		return nil, &services.RpcException{
			Code: kErrorCodeNotFound,
			Msg:  fmt.Sprintf("Lookup failed, Msg: %s", err.Error()),
		}
	} else {
		geoData := &GeoData{
			CountryName:    city.Country.Names["en"],
			CountryIsoCode: city.Country.IsoCode,
			CityName:       city.City.Names["en"],
			Lat:            fmt.Sprintf("%.7f", city.Location.Latitude),
			Lng:            fmt.Sprintf("%.7f", city.Location.Longitude),
			Timezone:       city.Location.TimeZone,
			Continent:      city.Continent.Names["en"],
			ContinentCode:  city.Continent.Code,
		}
		if len(city.Subdivisions) > 0 {
			geoData.Province = city.Subdivisions[0].IsoCode
		}

		data, _ := json.Marshal(geoData)
		log.Debugf("IP: %s, result found: %s, elapsed: %.3fms", ipStr, string(data), float64(microseconds() - start) * 0.001)

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
