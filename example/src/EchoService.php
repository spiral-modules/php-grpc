<?php
/**
 * Sample GRPC PHP server.
 */

use Spiral\GRPC\ContextInterface;

class EchoService implements Service\EchoInterface
{
    public function Ping(ContextInterface $ctx, Message $in): Message
    {
        error_log(print_r($ctx, 1));

        return $in;
    }
}