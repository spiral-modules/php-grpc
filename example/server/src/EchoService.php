<?php
/**
 * Sample GRPC PHP server.
 */

use Spiral\GRPC\ContextInterface;
use Service\Message;

class EchoService implements Service\EchoInterface
{
    public function Ping(ContextInterface $ctx, Message $in): Message
    {
        $out = new Message();
        return $out->setMsg(strrev($in->getMsg()));
    }
}