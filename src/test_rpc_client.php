<?php

require(__DIR__ . '/vendor/autoload.php');

use Thrift\Exception\TException;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Protocol\TMultiplexedProtocol;
use Thrift\Transport\TSocket;
use Thrift\Transport\TFramedTransport;
use Thrift\ClassLoader\ThriftClassLoader;

$loader = new ThriftClassLoader();
$loader->registerDefinition('geoip_service', __DIR__ . '/gen-php/');
$loader->registerDefinition('rpc_thrift', __DIR__ . '/gen-php/');

// 注册loader
$loader->register();

use geoip_service\GeoIpServiceClient;


class TestCode {

  function testProxiedRPCHelloworld() {
    try {

      // 直接使用rpc proxy进行通信
      $socket = new TSocket('tcp://localhost', 5550);
      // $socket = new TSocket('/usr/local/rpc_proxy/proxy.sock');

      $transport = new TFramedTransport($socket, true, true);

      // 指定后端服务
      $service_name = "ipgeo";
      $protocol = new TMultiplexedProtocol(new TBinaryProtocol($transport), $service_name);

      // 创建Client
      $client = new GeoIpServiceClient($protocol);

      $transport->open();
      $data = $client->IpToGeoData("218.97.243.4");
      var_dump($data);

      $transport->close();

    } catch (TException $tx) {
      print 'TException: ' . $tx->getMessage() . "\n";
    }
  }
}

$test_code = new TestCode();
$test_code->testProxiedRPCHelloworld();
