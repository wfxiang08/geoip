package geoip

import (
	"encoding/json"
	"fmt"
	geoip_services "geoip/services"
	"github.com/fatih/color"
	"github.com/getsentry/raven-go"
	"github.com/oschwald/maxminddb-golang"
	"github.com/wfxiang08/cyutils/utils/atomic2"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"github.com/wfxiang08/thrift_rpc_base/rpcthrift/services"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

// 只是为了方便定位
const (
	kErrorCodeNotFound = 1111
)

var reader *maxminddb.Reader
var rwLock sync.RWMutex

func microseconds() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

// 内部函数不要随便返回一般的Error, 统一封装成为: RpcException
// 注意陷阱: https://golang.org/doc/faq#nil_error
//
func IpToGeoData(ipStr string) (*geoip_services.GeoData, error) {
	start := microseconds()
	ip := net.ParseIP(ipStr)
	var city City

	rwLock.RLock()
	err := reader.Lookup(ip, &city)
	rwLock.RUnlock()

	if err != nil {
		log.Debugf("WARNGIN: IP: %s, no result found, elapsed: %.3fms", ipStr, float64(microseconds()-start)*0.001)
		return nil, &services.RpcException{
			Code: kErrorCodeNotFound,
			Msg:  fmt.Sprintf("Lookup failed, Msg: %s", err.Error()),
		}
	} else {
		geoData := &geoip_services.GeoData{
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
		log.Debugf("IP: %s, result found: %s, elapsed: %.3fms", ipStr, string(data), float64(microseconds()-start)*0.001)

		return geoData, nil
	}
}

// 全内存保存DB数据
func InitMaxMindDb(dbPath string, shouldPanic bool) {
	data, err := ioutil.ReadFile(dbPath)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		if shouldPanic {
			log.PanicError(err, "Read GeoIP failed")
		} else {
			log.ErrorErrorf(err, "Read GeoIP failed")
			return
		}
	}
	reader1, err1 := maxminddb.FromBytes(data)
	if err1 != nil {
		raven.CaptureErrorAndWait(err1, nil)
		if shouldPanic {
			log.PanicError(err1, "Read GeoIP failed")
		} else {
			log.ErrorErrorf(err1, "Read GeoIP failed")
			return
		}
	}

	raven.CaptureMessage("InitMaxMindDb succeed", nil)
	rwLock.Lock()
	reader = reader1
	rwLock.Unlock()
}

func StartHttpProfile(profileAddr string, running atomic2.Bool) {
	// 异步启动一个debug server
	// 由于可能多个进程同时存在，因此如果失败了，就等待重启
	log.Printf(color.RedString("Profile Address: %s"), profileAddr)
	for running.Get() {
		time.Sleep(time.Second * 10)
		err := http.ListenAndServe(profileAddr, nil)
		if err != nil {
			log.ErrorErrorf(err, "profile server start failed")
		} else {
			break
		}
	}
}
