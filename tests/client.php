<?php

use Spiral\Goridge;

ini_set('display_errors', 'stderr');
require dirname(__DIR__) . "/vendor/autoload.php";

require 'GPBMetadata/Test.php';
require 'Test/Message.php';

$relay = new Goridge\StreamRelay(STDIN, STDOUT);

$w = new \Spiral\RoadRunner\Worker($relay);

while ($body = $w->receive($context)) {
    try {
        $o = new \Test\Message();
        $o->setMsg("hi from php");

        $w->send($o->serializeToString());
    } catch (\Throwable $e) {
        $w->error((string)$w);
    }
}
