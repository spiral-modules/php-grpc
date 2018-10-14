<?php

ini_set('display_errors', 'stderr');
require "vendor/autoload.php";


class Yo
{
    public function ping(\Spiral\GRPC\Context $context, \Test\Message $msg): \Test\Message
    {
        return $msg;
    }

    public function ping2(\Test\Message $msg): \Test\Message
    {
        return $msg;
    }
}

$s = new \Spiral\GRPC\Service("yo", new Yo());


dump($s);