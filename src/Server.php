<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\GRPCException;
use Spiral\GRPC\Exception\NotFoundException;
use Spiral\GRPC\Exception\ServiceException;
use Spiral\RoadRunner\Worker;

/**
 * Manages group of services and communication with RoadRunner server.
 */
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
        $this->invocator = $invocator ?? new Invocator();
    }

    /**
     * Register new GRPC service.
     *
     * Example: $server->registerService(EchoServiceInterface::class, new EchoService());
     *
     * @param string $interface Generated service interface.
     * @param object $handler Must implement interface.
     *
     * @throws ServiceException
     */
    public function registerService(string $interface, $handler)
    {
        $service = new Service($this->invocator, $interface, $handler);
        $this->services[$service->getName()] = $service;
    }

    /**
     * Serve GRPC over given RoadRunner worker.
     *
     * @param Worker $worker
     */
    public function serve(Worker $worker)
    {
        while ($body = $worker->receive($ctx)) {
            try {
                $ctx = json_decode($ctx, true);
                $worker->send($this->invoke($ctx['service'], $ctx['method'], $ctx['context'], $body));
            } catch (GRPCException $e) {
                $worker->send($this->packError($e));
            } catch (\Throwable $e) {
                $worker->error($e);
            }
        }
    }

    /**
     * Invoke service method with binary payload and return the response.
     *
     * @param string $service
     * @param string $method
     * @param array|null $context
     * @param string $body
     * @return string
     *
     * @throws GRPCException
     * @throws \Throwable
     */
    protected function invoke(string $service, string $method, array $context = null, string $body): string
    {
        if (!isset($this->services[$service])) {
            throw new NotFoundException("Service `{$service}` not found.", StatusCode::NOT_FOUND);
        }

        return $this->services[$service]->invoke($method, new Context($context ?? []), $body);
    }

    /**
     * Packs exception message and code into one string.
     *
     * @param GRPCException $e
     * @return string
     */
    private function packError(GRPCException $e): string
    {
        return sprintf("%s|:|%s", $e->getCode(), $e->getMessage());
    }
}