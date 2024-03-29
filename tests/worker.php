<?php
/**
 * Sample GRPC PHP server.
 */

use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require __DIR__ . '/../vendor_php/autoload.php';

$server = new \Spiral\GRPC\Server();
$server->registerService(\Service\TestInterface::class, new \Test\TestService());
$server->registerService(\Health\HealthInterface::class, new \HealthImpl\HealthService());

$w = new RoadRunner\Worker(new Goridge\StreamRelay(STDIN, STDOUT));
$server->serve($w);
