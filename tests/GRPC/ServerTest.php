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
use Service\TestInterface;
use Spiral\GRPC\Server;
use Spiral\RoadRunner\Worker;
use Test\TestService;

class ServerTest extends TestCase
{
    public function testInvoke()
    {
        $s = new Server();
        $s->registerService(TestInterface::class, new TestService());

        $w = new TestWorker($this, [
            [
                'ctx'     => [
                    'service' => 'service.Test',
                    'method'  => 'Echo',
                    'context' => [],
                ],
                'send'    => $this->packMessage('hello world'),
                'receive' => $this->packMessage('hello world')
            ]
        ]);

        $s->serve($w);

        $this->assertTrue($w->done());
    }

    public function testNotFound()
    {
        $s = new Server();
        $s->registerService(TestInterface::class, new TestService());

        $w = new TestWorker($this, [
            [
                'ctx'   => [
                    'service' => 'service.Test2',
                    'method'  => 'Echo',
                    'context' => [],
                ],
                'send'  => $this->packMessage('hello world'),
                'error' => '5|:|Service `service.Test2` not found.'
            ]
        ]);

        $s->serve($w);

        $this->assertTrue($w->done());
    }

    public function testNotFound2()
    {
        $s = new Server();
        $s->registerService(TestInterface::class, new TestService());

        $w = new TestWorker($this, [
            [
                'ctx'   => [
                    'service' => 'service.Test',
                    'method'  => 'Echo2',
                    'context' => [],
                ],
                'send'  => $this->packMessage('hello world'),
                'error' => '5|:|Method `Echo2` not found in service `service.Test`.'
            ]
        ]);

        $s->serve($w);

        $this->assertTrue($w->done());
    }

    private function packMessage(string $message): string
    {
        $m = new Message();
        $m->setMsg($message);

        return $m->serializeToString();
    }
}

class TestWorker extends Worker
{
    private $t;
    private $sequence = [];
    private $pos = 0;

    public function __construct(TestCase $t, array $sequence)
    {
        $this->t = $t;
        $this->sequence = $sequence;
    }

    public function done()
    {
        return $this->pos == count($this->sequence);
    }

    public function receive(&$header)
    {
        if (!isset($this->sequence[$this->pos])) {
            $header = null;
            return null;
        }

        $header = json_encode($this->sequence[$this->pos]['ctx']);

        return $this->sequence[$this->pos]['send'];
    }

    public function send(string $payload = null, string $header = null)
    {
        $this->t->assertSame($this->sequence[$this->pos]['receive'], $payload);
        $this->pos++;
    }

    public function error(string $message)
    {
        $this->t->assertSame($this->sequence[$this->pos]['error'], $message);
        $this->pos++;
    }
}