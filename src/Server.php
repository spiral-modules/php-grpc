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
    /** @var InvokerInterface */
    private $invocator;

    /** @var ServiceWrapper[] */
    private $services = [];

    /**
     * @param InvokerInterface|null $invocator
     */
    public function __construct(InvokerInterface $invocator = null)
    {
        $this->invocator = $invocator ?? new Invoker();
    }

    /**
     * Register new GRPC service.
     *
     * Example: $server->registerService(EchoServiceInterface::class, new EchoService());
     *
     * @param string           $interface Generated service interface.
     * @param ServiceInterface $service   Must implement interface.
     *
     * @throws ServiceException
     */
    public function registerService(string $interface, ServiceInterface $service)
    {
        $service = new ServiceWrapper($this->invocator, $interface, $service);
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
                $worker->send($this->invoke($ctx['service'], $ctx['method'], $ctx['context'],
                    $body));
            } catch (GRPCException $e) {
                $worker->error($this->packError($e));
            } catch (\Throwable $e) {
                $worker->error($e);
            }
        }
    }

    /**
     * Invoke service method with binary payload and return the response.
     *
     * @param string     $service
     * @param string     $method
     * @param array|null $context
     * @param string     $body
     * @return string
     *
     * @throws GRPCException
     * @throws \Throwable
     */
    protected function invoke(
        string $service,
        string $method,
        array $context = null,
        string $body
    ): string {
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