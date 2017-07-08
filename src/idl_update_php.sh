#!/usr/bin/env bash
thrift -r --gen php Geoip.Services.thrift
rm -rf ../Geoip
mv gen-php/Geoip ../Geoip
# rm -rf gen-php

