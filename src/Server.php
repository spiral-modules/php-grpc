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
    private $invoker;

    /** @var ServiceWrapper[] */
    private $services = [];

    /**
     * @param InvokerInterface|null $invoker
     */
    public function __construct(InvokerInterface $invoker = null)
    {
        $this->invoker = $invoker ?? new Invoker();
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
        $service = new ServiceWrapper($this->invoker, $interface, $service);
        $this->services[$service->getName()] = $service;
    }

    /**
     * Serve GRPC over given RoadRunner worker.
     *
     * @param Worker $worker
     */
    public function serve(Worker $worker)
    {
        while (true) {
            $body = $worker->receive($ctx);
            if (empty($body) && empty($ctx)) {
                return;
            }

            try {
                $ctx = json_decode($ctx, true);
                $resp = $this->invoke(
                    $ctx['service'],
                    $ctx['method'],
                    $ctx['context'] ?? [],
                    $body
                );

                $worker->send($resp);
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
     * @param string $service
     * @param string $method
     * @param array  $context
     * @param string $body
     * @return string
     *
     * @throws GRPCException
     * @throws \Throwable
     */
    protected function invoke(
        string $service,
        string $method,
        array $context,
        ?string $body
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