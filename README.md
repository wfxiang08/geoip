#geoip rpc
* 基于GeoIP2-City.mmdb实现从ip到geo信息的映射
* 服务器采用golang开发，并且通过rpc对外提供服务

## Geoip服务运行:
* cd geoip/src
* source start_env.sh
* glide install
* go build cmd/service_geoip.go
* 运行:
	* 配合rpc_proxy框架运行，会自动将服务注册到zk中; 为长连接模式下的请求而优化
		* ./service_geoip -c config-service.ini
    * 独立运行，不依赖zookeeper
		* ./service_geoip -c config-standalone.ini
* Php版本的Client:
	* cd geoip
	* composer install
	* php test_rpc_client.php
	* 注意: autoload.php的使用
