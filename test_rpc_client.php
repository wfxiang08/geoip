<?php

require(__DIR__.'/vendor/autoload.php');

use Geoip\Services\GeoIpServiceClient;
use Thrift\Exception\TException;
use Thrift\Protocol\TBinaryProtocolAccelerated;

class TestCode {
  function testDirectGeoIP() {
    $this->t1();
    $this->t2();

    echo "================\n";
    $this->t1();
    $this->t2();
  }

  function t1() {
    echo "testing t1\n";
    try {
      $iteration_num = 100;
      $start1 = microtime(true);
      $sock = new SmSocket('127.0.0.1', 5563, true, true);
      $sock->pconnect(200);
      $start1 = microtime(true) - $start1;

      $start2 = microtime(true);
      $sock = new SmSocket('127.0.0.1', 5563, true, true);
      $sock->pconnect(200);
      $start2 = microtime(true) - $start2;


      $start = microtime(true);
      $client = new GeoIpServiceClient($sock);

      for ($i = 0; $i < 100; $i++) {
        $data = $client->IpToGeoData("218.97.243.4");
      }

      $start = microtime(true) - $start;
      echo "Open: ".sprintf("%.3fms", $start1 * 1000).", Reopen: ".sprintf("%.3fms", $start2 * 1000).", IpToGeoData: ".sprintf("%.3fms\n", $start / $iteration_num * 1000);

    } catch (TException $tx) {
      print 'TException: '.$tx->getMessage()."\n";
    }
  }

  function t2() {
    echo "testing t2\n";
    try {
      $iteration_num = 100;
      $start1 = microtime(true);
      $sock = new \Thrift\Transport\TSocket('tcp://127.0.0.1', 5563);
      $sock->open();
      $start1 = microtime(true) - $start1;

      $start2 = microtime(true);
      $sock = new \Thrift\Transport\TSocket('tcp://127.0.0.1', 5563);
      $sock->open();
      $start2 = microtime(true) - $start2;


      $start = microtime(true);
      $transport = new \Thrift\Transport\TFramedTransport($sock);
      $protocol = new \Thrift\Protocol\TBinaryProtocol($transport);
      $client = new GeoIpServiceClient($protocol);

      for ($i = 0; $i < 100; $i++) {
        $data = $client->IpToGeoData("218.97.243.4");
      }

      $start = microtime(true) - $start;
      echo "OpenClose: ".sprintf("%.3fms", $start1 * 1000).", Reopen: ".sprintf("%.3fms", $start2 * 1000).", IpToGeoData: ".sprintf("%.3fms\n", $start / $iteration_num * 1000);

    } catch (TException $tx) {
      print 'TException: '.$tx->getMessage()."\n";
    }
  }
}

$test_code = new TestCode();
$test_code->testDirectGeoIP();
