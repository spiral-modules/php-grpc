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
        // TODO: Implement Die() method.
    }

    public function Info(ContextInterface $ctx, Message $in): Message
    {
        // TODO: Implement Info() method.
    }
}