<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC\Tests;

use PHPUnit\Framework\TestCase;
use Service\Message;
use Spiral\GRPC\Context;
use Spiral\GRPC\Invoker;
use Spiral\GRPC\Method;
use Test\TestService;

class InvokerTest extends TestCase
{
    public function testInvoke(): void
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Echo'));

        $i = new Invoker();

        $out = $i->invoke($s, $m, new Context([]), $this->packMessage('hello'));

        $m = new Message();
        $m->mergeFromString($out);

        $this->assertSame('hello', $m->getMsg());
    }

    public function testInvokeError(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\InvokeException::class);

        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Echo'));

        $i = new Invoker();

        $i->invoke($s, $m, new Context([]), 'invalid-message');
    }

    private function packMessage(string $message): string
    {
        $m = new Message();
        $m->setMsg($message);

        return $m->serializeToString();
    }
}
