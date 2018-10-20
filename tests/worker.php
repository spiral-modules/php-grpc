<?php
/**
 * Sample GRPC PHP server.
 */

use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require "../vendor/autoload.php";

$server = new \Spiral\GRPC\Server();
$server->registerService(\Service\TestService::class, new TestService());

$w = new RoadRunner\Worker(new Goridge\StreamRelay(STDIN, STDOUT));
$server->serve($w);