package main

import (
	ips "geoip/services"
	"github.com/wfxiang08/rpc_proxy/src/proxy"
	"geoip"
	"flag"
)

const (
	BINARY_NAME = "service_geoip"
	SERVICE_DESC = "GeoIp Service v0.1"
	IP_DATA = "/usr/local/ip/GeoIP2-City.mmdb"
)

var (
	buildDate string
	gitVersion string
	ipData = flag.String("ip", IP_DATA, "Ip Data file")
)

func main() {

	proxy.RpcMain(BINARY_NAME, SERVICE_DESC,
		// 默认的ThriftServer的配置checker
		proxy.ConfigCheckThriftService,

		// 可以根据配置config来创建processor
		func(config *proxy.ServiceConfig) proxy.Server {
			handler := geoip.NewHandler(*ipData)
			processor := ips.NewGeoIpServiceProcessor(handler)
			return proxy.NewThriftRpcServer(config, processor)
		}, buildDate, gitVersion)
}

