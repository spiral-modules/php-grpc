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
use Service\TestInterface;
use Spiral\GRPC\Context;
use Spiral\GRPC\Invoker;
use Spiral\GRPC\ServiceInterface;
use Spiral\GRPC\ServiceWrapper;
use Test\TestService;

class ServiceWrapperTest extends TestCase implements ServiceInterface
{
    public function testName(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            new TestService()
        );

        $this->assertSame('service.Test', $w->getName());
    }

    public function testService(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            $t = new TestService()
        );

        $this->assertSame($t, $w->getService());
    }

    public function testMethods(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            new TestService()
        );

        $this->assertCount(5, $w->getMethods());
    }

    public function testInvokeNotFound(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\NotFoundException::class);

        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            new TestService()
        );

        $w->invoke('NotFound', new Context([]), '');
    }

    public function testInvoke(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            new TestService()
        );

        $out = $w->invoke('Echo', new Context([]), $this->packMessage('hello world'));

        $m = new Message();
        $m->mergeFromString($out);

        $this->assertSame('hello world', $m->getMsg());
    }

    public function testNotImplemented(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\ServiceException::class);

        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            $this
        );
    }

    public function testInvalidInterface(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\ServiceException::class);

        $w = new ServiceWrapper(
            new Invoker(),
            InvalidInterface::class,
            $this
        );
    }

    public function testInvalidInterface2(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\ServiceException::class);

        $w = new ServiceWrapper(
            new Invoker(),
            'NotFound',
            $this
        );
    }

    private function packMessage(string $message): string
    {
        $m = new Message();
        $m->setMsg($message);

        return $m->serializeToString();
    }
}
