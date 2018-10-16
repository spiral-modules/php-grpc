<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Spiral\RoadRunner\Worker;

class Server
{
    /** @var InvocatorInterface */
    private $invocator;

    /** @var Service[] */
    private $services = [];

    /**
     * @param InvocatorInterface|null $invocator
     */
    public function __construct(InvocatorInterface $invocator = null)
    {
        $this->invocator = $invocator;
    }

    /**
     * @param string $name
     * @param object $handler
     */
    public function registerService(string $name, $handler)
    {
        $this->services[$name] = new Service($name, $handler, $this->invocator);
    }

    public function serve(Worker $worker)
    {
        while ($body = $worker->receive($context)) {
            try {
                $context = json_decode($context, true);

                $worker->send($this->invoke(
                    $context['service'],
                    $context['method'],
                    $body,
                    $context['context']
                ));
            } catch (\Throwable $e) {
                // report error
                // todo: map error to GRPC errors
            }
        }
    }

    protected function invoke(string $service, string $method, string $body, array $context): string
    {
        //todo: check if service exists
        return $this->services[$service]->invoke($method, new Context($context), $body);
    }
}