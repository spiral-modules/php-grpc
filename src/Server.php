<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Google\Protobuf\Any;
use Spiral\GRPC\Exception\GRPCException;
use Spiral\GRPC\Exception\GRPCExceptionInterface;
use Spiral\GRPC\Exception\NotFoundException;
use Spiral\GRPC\Exception\ServiceException;
use Spiral\GRPC\Internal\Json;
use Spiral\RoadRunner\Payload;
use Spiral\RoadRunner\Worker;

/**
 * Manages group of services and communication with RoadRunner server.
 *
 * @psalm-type ServerOptions = array {
 *  debug?: bool
 * }
 *
 * @psalm-type ContextResponse = array {
 *  service: string,
 *  method:  string,
 *  context: array<string, array<string>>
 * }
 */
final class Server
{
    /**
     * @var InvokerInterface
     */
    private $invoker;

    /**
     * @var ServiceWrapper[]
     */
    private $services = [];

    /**
     * @var ServerOptions
     */
    private $options;

    /**
     * @param InvokerInterface|null $invoker
     * @param ServerOptions $options
     */
    public function __construct(InvokerInterface $invoker = null, array $options = [])
    {
        $this->invoker = $invoker ?? new Invoker();
        $this->options = $options;
    }

    /**
     * Register new GRPC service.
     *
     * For example:
     * <code>
     *  $server->registerService(EchoServiceInterface::class, new EchoService());
     * </code>
     *
     * @template T of ServiceInterface
     *
     * @param class-string<T> $interface Generated service interface.
     * @param T $service Must implement interface.
     * @throws ServiceException
     */
    public function registerService(string $interface, ServiceInterface $service): void
    {
        $service = new ServiceWrapper($this->invoker, $interface, $service);

        $this->services[$service->getName()] = $service;
    }

    /**
     * @param string $body
     * @param ContextResponse $data
     * @return array{ 0: string, 1: string }
     * @throws \JsonException
     * @throws \Throwable
     */
    private function tick(string $body, array $data): array
    {
        $context = (new Context($data['context']))
            ->withValue(ResponseHeaders::class, new ResponseHeaders())
        ;

        $response = $this->invoke($data['service'], $data['method'], $context, $body);

        /** @var ResponseHeaders|null $responseHeaders */
        $responseHeaders = $context->getValue(ResponseHeaders::class);
        $responseHeadersString = $responseHeaders ? $responseHeaders->packHeaders() : '{}';

        return [$response, $responseHeadersString];
    }

    /**
     * @param Worker $worker
     * @param string $body
     * @param string $headers
     * @psalm-suppress InaccessibleMethod
     */
    private function workerSend(Worker $worker, string $body, string $headers): void
    {
        // RoadRunner 1.x
        if (\method_exists($worker, 'send')) {
            $worker->send($body, $headers);

            return;
        }

        // RoadRunner 2.x
        $worker->respond(new Payload($body, $headers));
    }

    /**
     * @param Worker $worker
     * @param string $message
     */
    private function workerError(Worker $worker, string $message): void
    {
        $worker->error($message);
    }

    /**
     * @param Worker $worker
     * @return array { 0: string, 1: string } | null
     *
     * @psalm-suppress UndefinedMethod
     * @psalm-suppress PossiblyUndefinedVariable
     */
    private function workerReceive(Worker $worker): ?array
    {
        /** @var string|\Stringable $body */
        $body = $worker->receive($ctx);

        if (empty($body) && empty($ctx)) {
            return null;
        }

        return [(string)$body, (string)$ctx];
    }

    /**
     * Serve GRPC over given RoadRunner worker.
     *
     * @param Worker $worker
     * @param callable|null $finalize
     */
    public function serve(Worker $worker, callable $finalize = null): void
    {
        while (true) {
            $request = $this->workerReceive($worker);

            if (! $request) {
                return;
            }

            [$body, $headers] = $request;

            try {
                /** @var ContextResponse $context */
                $context = Json::decode((string)$headers);

                [$answerBody, $answerHeaders] = $this->tick((string)$body, $context);

                $this->workerSend($worker, $answerBody, $answerHeaders);
            } catch (GRPCExceptionInterface $e) {
                $this->workerError($worker, $this->packError($e));
            } catch (\Throwable $e) {
                $this->workerError($worker, $this->isDebugMode() ? (string)$e : $e->getMessage());
            } finally {
                if ($finalize !== null) {
                    isset($e) ? $finalize($e) : $finalize();
                }
            }
        }
    }

    /**
     * Invoke service method with binary payload and return the response.
     *
     * @param string $service
     * @param string $method
     * @param ContextInterface $context
     * @param string $body
     * @return string
     * @throws GRPCException
     */
    protected function invoke(string $service, string $method, ContextInterface $context, string $body): string
    {
        if (! isset($this->services[$service])) {
            throw NotFoundException::create("Service `{$service}` not found.", StatusCode::NOT_FOUND);
        }

        return $this->services[$service]->invoke($method, $context, $body);
    }

    /**
     * Packs exception message and code into one string.
     *
     * Internal agreement:
     *
     * Details will be sent as serialized google.protobuf.Any messages after
     * code and exception message separated with |:| delimiter.
     *
     * @param GRPCExceptionInterface $e
     * @return string
     */
    private function packError(GRPCExceptionInterface $e): string
    {
        $data = [$e->getCode(), $e->getMessage()];

        foreach ($e->getDetails() as $detail) {
            $anyMessage = new Any();

            $anyMessage->pack($detail);

            $data[] = $anyMessage->serializeToString();
        }

        return \implode('|:|', $data);
    }

    /**
     * If server runs in debug mode
     *
     * @return bool
     */
    private function isDebugMode(): bool
    {
        $debug = false;

        if (isset($this->options['debug'])) {
            $debug = \filter_var($this->options['debug'], \FILTER_VALIDATE_BOOLEAN);
        }

        return $debug;
    }
}
