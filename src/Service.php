<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

class Service
{
    private $handler;

    public function __construct($handler)
    {
        $this->handler = $handler;
    }


}