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
      // XXX:
      $service_name = "ipgeo";
      $protocol = new TMultiplexedProtocol(new TBinaryProtocol($transport), $service_name);

      // 创建Client
      // XXX:
      $client = new GeoIpServiceClient($protocol);

      $transport->open();

      // 同步调用
      $data = $client->IpToGeoData("218.97.243.4");
      var_dump($data);

      // 异步调用
      echo "Send request1 async\n";
      $client->send_IpToGeoData("218.97.243.4");
      echo "Send request2 async\n";
      $client->send_IpToGeoData("106.201.41.11");

      echo "Receive request1 async\n";
      $data1 = $client->recv_IpToGeoData();
      echo "Receive request2 async\n";
      $data2 = $client->recv_IpToGeoData();
      var_dump($data1);
      var_dump($data2);

      $transport->close();

    } catch (TException $tx) {
      print 'TException: ' . $tx->getMessage() . "\n";
    }
  }
}

$test_code = new TestCode();
$test_code->testProxiedRPCHelloworld();
