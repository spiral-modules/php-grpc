<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC\Tests;

use PHPUnit\Framework\TestCase;
use Spiral\GRPC\Context;

class ContextTest extends TestCase
{
    public function testGetValue()
    {
        $ctx = new Context([
            'key' => ['value']
        ]);

        $this->assertSame(['value'], $ctx->getValue('key'));
    }

    public function testGetNullValue()
    {
        $ctx = new Context([
            'key' => ['value']
        ]);

        $this->assertSame(null, $ctx->getValue('other'));
    }

    public function testGetValues()
    {
        $ctx = new Context([
            'key' => ['value']
        ]);

        $this->assertSame([
            'key' => ['value']
        ], $ctx->getValues());
    }


    public function testWithValue()
    {
        $ctx = new Context([
            'key' => ['value']
        ]);

        $this->assertSame(['value'], $ctx->getValue('key'));

        $ctx2 = $ctx->withValue('new', 'another')->withValue('key', ['value2']);

        $this->assertSame(['value'], $ctx->getValue('key'));
        $this->assertSame(null, $ctx->getValue('new'));

        $this->assertSame(['value2'], $ctx2->getValue('key'));
        $this->assertSame('another', $ctx2->getValue('new'));
    }
}