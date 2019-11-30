<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Google\Protobuf\Any;
use Google\Protobuf\Internal\Message;
use Spiral\GRPC\Exception\GRPCException;
use Spiral\GRPC\Exception\NotFoundException;
use Spiral\GRPC\Exception\ServiceException;
use Spiral\RoadRunner\Worker;

/**
 * Manages group of services and communication with RoadRunner server.
 */
final class Server
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
     * @param ServiceInterface $service Must implement interface.
     *
     * @throws ServiceException
     */
    public function registerService(string $interface, ServiceInterface $service): void
    {
        $service = new ServiceWrapper($this->invoker, $interface, $service);
        $this->services[$service->getName()] = $service;
    }

    /**
     * Serve GRPC over given RoadRunner worker.
     *
     * @param Worker        $worker
     * @param callable|null $finalize
     */
    public function serve(Worker $worker, callable $finalize = null): void
    {
        while (true) {
            $body = $worker->receive($ctx);
            if (empty($body) && empty($ctx)) {
                return;
            }

            try {
                $ctx = json_decode($ctx, true);
                $grpcCtx = new Context(
                    $ctx['context'] + [ResponseHeaders::class => new ResponseHeaders([])]
                );

                $resp = $this->invoke(
                    $ctx['service'],
                    $ctx['method'],
                    $grpcCtx,
                    $body
                );

                /** @var ResponseHeaders|null $responseHeaders */
                $responseHeaders = $grpcCtx->getValue(ResponseHeaders::class);
                $worker->send($resp, $responseHeaders ? $responseHeaders->packHeaders() : '{}');
            } catch (GRPCException $e) {
                $worker->error($this->packError($e));
            } catch (\Throwable $e) {
                $worker->error((string)$e);
            } finally {
                if ($finalize !== null) {
                    call_user_func($finalize, $e ?? null);
                }
            }
        }
    }

    /**
     * Invoke service method with binary payload and return the response.
     *
     * @param string  $service
     * @param string  $method
     * @param Context $context
     * @param string  $body
     * @return string
     *
     * @throws GRPCException
     * @throws \Throwable
     */
    protected function invoke(
        string $service,
        string $method,
        Context $context,
        ?string $body
    ): string {
        if (!isset($this->services[$service])) {
            throw new NotFoundException("Service `{$service}` not found.", StatusCode::NOT_FOUND);
        }

        return $this->services[$service]->invoke($method, $context, $body);
    }

    /**
     * Packs exception message and code into one string.
     *
     * Internal agreement:
     *
     * Details will be sent as serialized google.protobuf.Any messages after code and exception message
     * separated with |:| delimeter.
     *
     * @param GRPCException $e
     * @return string
     */
    private function packError(GRPCException $e): string
    {
        $data = [$e->getCode(), $e->getMessage()];

        foreach ($e->getDetails() as $detail) {
            /**
             * @var Message $detail
             */
            $anyMessage = new Any();

            $anyMessage->pack($detail);

            $data[] = $anyMessage->serializeToString();
        }

        return implode('|:|', $data);
    }
}
