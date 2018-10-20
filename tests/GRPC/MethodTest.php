<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC\Tests;

use PHPUnit\Framework\TestCase;
use Service\Message;
use Spiral\GRPC\ContextInterface;
use Spiral\GRPC\Method;
use Test\TestService;

class MethodTest extends TestCase
{
    /**
     * @expectedException \Spiral\GRPC\Exception\GRPCException
     */
    public function testInvalidParse()
    {
        Method::parse(new \ReflectionMethod($this, 'testInvalidParse'));
    }

    public function testMatch()
    {
        $s = new TestService();
        $this->assertTrue(Method::match(new \ReflectionMethod($s, 'Info')));
    }

    public function testNoMatch()
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 't_M')));
    }

    public function testNoMatch2()
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 't_M2')));
    }

    public function testNoMatch3()
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 't_M3')));
    }

    public function testNoMatch4()
    {
        $this->assertFalse(Method::match(new \ReflectionMethod($this, 't_M4')));
    }

    public function testMethodName()
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame('Info', $m->getName());
    }

    public function testMethodInputType()
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame(Message::class, $m->getInputType());
    }

    public function testMethodOutputType()
    {
        $s = new TestService();
        $m = Method::parse(new \ReflectionMethod($s, 'Info'));
        $this->assertSame(Message::class, $m->getOutputType());
    }

    public function t_M(ContextInterface $context, TestService $input): Message
    {

    }

    public function t_M2(ContextInterface $context, Message $input): TestService
    {

    }

    public function t_M3(TestService $context, Message $input): TestService
    {

    }

    public function t_M4(TestService $context, Message $input): Invalid
    {

    }
}