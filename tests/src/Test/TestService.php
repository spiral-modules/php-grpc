<?php

namespace Test;

use Service\DetailsMessageForException;
use Service\EmptyMessage;
use Service\Message;
use Service\TestInterface;
use Spiral\GRPC;
use Spiral\GRPC\ContextInterface;
use Spiral\GRPC\Exception\GRPCException;
use Spiral\GRPC\Exception\NotFoundException;

class TestService implements TestInterface
{
    public function Echo(ContextInterface $ctx, Message $in, array &$metadata = []): Message
    {
        return $in;
    }

    public function Throw(ContextInterface $ctx, Message $in, array &$metadata = []): Message
    {
        $out = new Message();

        switch ($in->getMsg()) {
            case "notFound":
                throw new NotFoundException("nothing here");
            case "withDetails":
                $detailsMessage = new DetailsMessageForException();
                $detailsMessage->setCode(1);
                $detailsMessage->setMessage("details message");

                $grpcException = new GRPCException("main exception message", 3, [$detailsMessage]);

                throw $grpcException;
        }

        return $out;
    }

    public function Die(ContextInterface $ctx, Message $in, array &$metadata = []): Message
    {
        error_log($in->getMsg());

        return $in;
    }

    public function Info(ContextInterface $ctx, Message $in, array &$metadata = []): Message
    {
        $out = new Message();
        switch ($in->getMsg()) {
            case "RR_GRPC":
            case "ENV_KEY":
                $out->setMsg(getenv($in->getMsg()));
                break;
            case "PID":
                $out->setMsg(getmypid());
                break;
            case"MD":
                $out->setMsg(json_encode($ctx->getValue('key')));
                break;
        }

        $metadata['foo'] = 'bar';

        return $out;
    }

    public function Ping(GRPC\ContextInterface $ctx, EmptyMessage $in, array &$metadata = []): EmptyMessage
    {
        return new EmptyMessage();
    }
}