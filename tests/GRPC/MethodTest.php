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
use Spiral\GRPC\ContextInterface;
use Spiral\GRPC\Method;
use Test\TestService;

class MethodTest extends TestCase
{
    public function testInvalidParse(): void
    {
        $this->expectException(\Spiral\GRPC\Exception\GRPCException::class);

        Method::parse(new \ReflectionMethod($this, 'testInvalidParse'));
    }

    public function testMatch(): void
    {
        $s = new TestService();
        $this->assertTrue(Method::match(new \ReflectionMethod($s, 'Info')));
    }

    public function testNoMatch(): void
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 'tM')));
    }

    public function testNoMatch2(): void
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 'tM2')));
    }

    public function testNoMatch3(): void
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 'tM3')));
    }

    public function testNoMatch4(): void
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 'tM4')));
    }

    public function testNoMatch5(): void
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 'tM5')));
    }

    public function testMethodName(): void
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame('Info', $m->getName());
    }

    public function testMethodInputType(): void
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame(Message::class, $m->getInputType());
    }

    public function testMethodOutputType(): void
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame(Message::class, $m->getOutputType());
    }

    public function tM(ContextInterface $context, TestService $input): Message
    {
    }

    public function tM2(ContextInterface $context, Message $input): TestService
    {
    }

    public function tM3(TestService $context, Message $input): TestService
    {
    }

    public function tM4(TestService $context, Message $input): Invalid
    {
    }

    public function tM5(TestService $context, Message $input): void
    {
    }
}
