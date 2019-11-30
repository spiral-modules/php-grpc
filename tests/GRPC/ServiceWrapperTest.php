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

    /**
     * @expectedException \Spiral\GRPC\Exception\NotFoundException
     */
    public function testInvokeNotFound(): void
    {
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

    /**
     * @expectedException \Spiral\GRPC\Exception\ServiceException
     */
    public function testNotImplemented(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            TestInterface::class,
            $this
        );
    }

    /**
     * @expectedException \Spiral\GRPC\Exception\ServiceException
     */
    public function testInvalidInterface(): void
    {
        $w = new ServiceWrapper(
            new Invoker(),
            InvalidInterface::class,
            $this
        );
    }

    /**
     * @expectedException \Spiral\GRPC\Exception\ServiceException
     */
    public function testInvalidInterface2(): void
    {
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
