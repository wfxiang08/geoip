package logserver

import (
	"github.com/oschwald/maxminddb-golang"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"io/ioutil"
	"net"
	"strings"
)

var reader *maxminddb.Reader

func ComputerCountry(item *LogItem) {
	computedCountry := LookupISOCodeByIp(item.ClientIp)
	if len(computedCountry) == 0 {
		computedCountry = item.Country
	}
	computedCountry = strings.ToUpper(computedCountry)
	item.ComputedCountry = computedCountry
}

func LookupISOCodeByIp(ipStr string) string {
	ip := net.ParseIP(ipStr)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}
	err := reader.Lookup(ip, &record)
	if err != nil {
		return ""
	} else {
		return record.Country.ISOCode
	}
}

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
