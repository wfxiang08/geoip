#!/usr/bin/env bash
rm -rf gen-go
rm -rf geoip/services
thrift --gen go:package_prefix="github.com/wfxiang08/thrift_rpc_base/",thrift_import="github.com/wfxiang08/go_thrift/thrift" Geoip.Services.thrift

mv gen-go/geoip/services geoip/services
rm -rf gen-go