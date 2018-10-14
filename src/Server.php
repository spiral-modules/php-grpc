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
    /** @var Service[] */
    private $services = [];

    public function addService(string $name, $handler)
    {
        // todo: validate
        $this->services[$name] = new Service($name, $handler);
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
            }
        }
    }

    protected function invoke(string $service, string $method, string $body, array $context): string
    {
        return $this->services[$service]->invoke($method, new Context($context), $body);
    }
}