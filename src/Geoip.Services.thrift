namespace php Geoip.Services

include "RpcThrift.Services.thrift"
const string VERSION = "0.0.1"

/**
 * 输入和输出的结果
 */
struct GeoData {
	1:string country_name,
	2:string country_iso_code,
	3:string city_name,
	4:string lat,
	5:string lng,
	6:string timezone,
	7:string continent,
	8:string continent_code,
	9:string province,
}

struct LatLng {
	1:string lat,
	2:string lng,
}

service GeoIpService extends RpcThrift.Services.RpcServiceBase {
    GeoData IpToGeoData(1: string ip) throws (1: RpcThrift.Services.RpcException re),
    LatLng GetLatlng(1: string ip) throws (1: RpcThrift.Services.RpcException re),

    string GetCityName(1: string ip) throws (1: RpcThrift.Services.RpcException re),
    string GetCountryName(1: string ip) throws (1: RpcThrift.Services.RpcException re),
    string GetCountryCode(1: string ip) throws (1: RpcThrift.Services.RpcException re),
    string GetProvince(1: string ip) throws (1: RpcThrift.Services.RpcException re),
}
