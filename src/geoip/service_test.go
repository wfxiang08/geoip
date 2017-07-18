package geoip

import (
	"testing"
	"fmt"
)

//
// go test geoip -v -run "TestIpGeo"
//
func TestIpGeo(t *testing.T) {
	dbPath := "/usr/local/ip/GeoIP2-City.mmdb"
	InitMaxMindDb(dbPath)

	func() {
		data, err := IpToGeoData("218.97.243.4")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Data: %s\n", data.String())
		}
	}()

	func() {
		data, err := IpToGeoData("67.220.91.30")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Data: %s\n", data.String())
		}
	}()
}
