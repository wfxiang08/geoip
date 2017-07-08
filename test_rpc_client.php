<?php

require(__DIR__ . '/vendor/autoload.php');

use Geoip\Services\GeoIpServiceClient;
use Thrift\Exception\TException;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Protocol\TMultiplexedProtocol;
use Thrift\Transport\TFramedTransport;
use Thrift\Transport\TSocket;

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

  function testDirectGeoIP() {
    try {
      $start = microtime(true);
      $iteration_num = 10000;
      for ($i = 0; $i < 10000; $i++) {
        $socket = new TSocket('/usr/local/video/geoip.sock');
        $transport = new TFramedTransport($socket, true, true);
        $protocol = new TBinaryProtocol($transport);
        $client = new GeoIpServiceClient($protocol);
        $transport->open();
        $data = $client->IpToGeoData("218.97.243.4");

        $transport->close();
      }
      $start = microtime(true) - $start;
      echo "OpenClose, Elapsed: " . sprintf("%.3fms\n", $start / $iteration_num * 1000);

      var_dump($data);

      $start = microtime(true);
      $iteration_num = 10000;
      $socket = new TSocket('/usr/local/video/geoip.sock');
      $transport = new TFramedTransport($socket, true, true);
      $protocol = new TBinaryProtocol($transport);
      $client = new GeoIpServiceClient($protocol);
      $transport->open();
      for ($i = 0; $i < 10000; $i++) {
        $data = $client->IpToGeoData("218.97.243.4");
      }
      $transport->close();
      $start = microtime(true) - $start;
      echo "No OpenClose, Elapsed: " . sprintf("%.3fms\n", $start / $iteration_num * 1000);
      var_dump($data);

    } catch (TException $tx) {
      print 'TException: ' . $tx->getMessage() . "\n";
    }
  }
}

$test_code = new TestCode();
$test_code->testDirectGeoIP();
