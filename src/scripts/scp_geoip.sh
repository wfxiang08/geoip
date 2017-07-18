#!/usr/bin/env bash
if [ "$#" -ne 1 ]; then
    echo "Please input hostname"
    exit -1
fi

host_name=$1

ssh root@${host_name} "mkdir -p /usr/local/services/geoip/"
ssh root@${host_name} "chown worker.worker -r /usr/local/services/geoip/"

scp config_standalone.ini  root@${host_name}:/usr/local/services/geoip/config_standalone.ini
ssh root@${host_name} "rm -f /usr/local/services/geoip/service_geoip"
scp service_geoip root@${host_name}:/usr/local/services/geoip/service_geoip
