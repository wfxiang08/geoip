package main

import (
	"flag"
	"geoip"
	ips "geoip/services"
	"github.com/getsentry/raven-go"
	"github.com/wfxiang08/cyutils/utils/atomic2"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"github.com/wfxiang08/overseer"
	"github.com/wfxiang08/rpc_proxy/src/proxy"
	"os"
	"os/signal"
	"syscall"
)

var (
	ipData      = flag.String("ip", "/usr/local/ip/GeoIP2-City.mmdb", "Ip Data file")
	logPrefix   = flag.String("log", "", "log file prefix")
	logLevel    = flag.String("level", "", "log level")
	rpcEndpoint = flag.String("end", "", "rpc endpoint")
	profileAddr = flag.String("profile_address", "", "profile address")
	pidfile     = flag.String("pidfile", "", "pidfile")
	sentry      = flag.String("sentry", "", "sentry address")
)

func main() {
	flag.Parse()
	raven.SetDSN(*sentry)

	// 1. 解析Log相关的配置
	if len(*logPrefix) > 0 {
		f, err := log.NewRollingFile(*logPrefix, 3)
		if err != nil {
			log.PanicErrorf(err, "open rolling log file failed: %s", *logPrefix)
		} else {
			// 不能放在子函数中
			defer f.Close()
			log.StdLog = log.New(f, "")
		}
	}

	// 默认是Debug模式
	log.SetLevel(log.LEVEL_DEBUG)
	log.SetFlags(log.Flags() | log.Lshortfile)

	// set log level
	if len(*logLevel) > 0 {
		proxy.SetLogLevel(*logLevel)
	}

	var addresses []string
	// 添加rcp的地址
	addresses = append(addresses, *rpcEndpoint)

	overseer.Run(overseer.Config{
		Program:   gracefulServer, // 执行的函数体
		Addresses: addresses,
		Debug:     false,
		Pidfile:   *pidfile,
	})
}

func gracefulServer(state overseer.State) {

	// 整体系统的状态
	var running atomic2.Bool
	running.Set(true)

	// 1. 开启DB的更新机制
	go func() {
		signals := make(chan os.Signal)
		// 信号意义参考: http://blog.csdn.net/Gpengtao/article/details/7879364
		signal.Notify(signals, syscall.SIGHUP)

		for running.Get() {

			<-signals

			// 重新加载配置文件
			log.Printf("===> Reload db: %s", *ipData)
			geoip.InitMaxMindDb(*ipData, false)
		}
	}()

	// 2. 开启Profile
	if len(*profileAddr) > 0 {
		go geoip.StartHttpProfile(*profileAddr, running)
	}

	// 3. 启动rpc服务
	handler := geoip.NewHandler(*ipData)
	processor := ips.NewGeoIpServiceProcessor(handler)

	overseer.GracefulRunWithListener(state.Listeners[0], processor)

}
