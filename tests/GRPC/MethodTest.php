<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC\Tests;

use PHPUnit\Framework\TestCase;
use Spiral\GRPC\Method;

class MethodTest extends TestCase
{
    /**
     * @expectedException \Spiral\GRPC\Exception\GRPCException
     */
    public function testInvalidParse()
    {
        $m = new \ReflectionMethod($this, 'testInvalidParse');
        Method::parse($m);
    }
}