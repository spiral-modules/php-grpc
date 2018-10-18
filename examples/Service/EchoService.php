<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Service;

use Spiral\GRPC\ContextInterface;

class EchoService implements EchoInterface
{
    public function Ping(ContextInterface $ctx, Message $in): Message
    {
        return $in;
    }
}