<?php

namespace Test;

use Service\Message;
use Service\TestInterface;
use Spiral\GRPC\ContextInterface;

class TestService implements TestInterface
{
    public function Echo(ContextInterface $ctx, Message $in): Message
    {
        return $in;
    }

    public function Throw(ContextInterface $ctx, Message $in): Message
    {
        // TODO: Implement Throw() method.
    }

    public function Die(ContextInterface $ctx, Message $in): Message
    {
        error_log($in->getMsg());

        return $in;
    }

    public function Info(ContextInterface $ctx, Message $in): Message
    {
        $out = new Message();
        switch ($in->getMsg()) {
            case "RR_GRPC":
            case "ENV_KEY":
                $out->setMsg(getenv($in->getMsg()));
                break;
        }

        return $out;
    }
}