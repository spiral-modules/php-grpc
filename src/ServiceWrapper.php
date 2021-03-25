<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\InvokeException;
use Spiral\GRPC\Exception\NotFoundException;
use Spiral\GRPC\Exception\ServiceException;

/**
 * Wraps handlers methods.
 */
final class ServiceWrapper
{
    /** @var string */
    private $name;

    /** @var ServiceInterface */
    private $service;

    /** @var InvokerInterface */
    private $invoker;

    /** @var Method[] */
    private $methods;

    /**
     * @param InvokerInterface $invoker
     * @param string $interface Service interface name.
     * @param ServiceInterface $service
     *
     * @throws ServiceException
     */
    public function __construct(
        InvokerInterface $invoker,
        string $interface,
        ServiceInterface $service
    ) {
        $this->invoker = $invoker;
        $this->configure($interface, $service);
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @return ServiceInterface
     */
    public function getService(): ServiceInterface
    {
        return $this->service;
    }

    /**
     * @return array
     */
    public function getMethods(): array
    {
        return array_values($this->methods);
    }

    /**
     * @param string $method
     * @param ContextInterface $context
     * @param string|null $input
     * @param array $metadata
     * @return string
     *
     * @throws NotFoundException
     * @throws InvokeException
     */
    public function invoke(string $method, ContextInterface $context, ?string $input, array &$metadata = []): string
    {
        if (! isset($this->methods[$method])) {
            throw NotFoundException::create("Method `{$method}` not found in service `{$this->name}`.");
        }

        return $this->invoker->invoke($this->service, $this->methods[$method], $context, $input, $metadata);
    }

    /**
     * Configure service name and methods.
     *
     * @param string $interface
     * @param ServiceInterface $service
     *
     * @throws ServiceException
     */
    protected function configure(string $interface, ServiceInterface $service): void
    {
        try {
            $r = new \ReflectionClass($interface);

            if (! $r->hasConstant('NAME')) {
                $message = "Invalid service interface `{$interface}`, constant `NAME` not found.";
                throw ServiceException::create($message);
            }

            $this->name = $r->getConstant('NAME');
        } catch (\ReflectionException $e) {
            $message = "Invalid service interface `{$interface}`.";
            throw ServiceException::create($message, StatusCode::INTERNAL, $e);
        }

        if (! $service instanceof $interface) {
            throw ServiceException::create("Service handler does not implement `{$interface}`.");
        }

        $this->service = $service;

        // list of all available methods and their object types
        $this->methods = $this->fetchMethods($service);
    }

    /**
     * @param ServiceInterface $service
     * @return array
     */
    protected function fetchMethods(ServiceInterface $service): array
    {
        $reflection = new \ReflectionObject($service);

        $methods = [];
        foreach ($reflection->getMethods(\ReflectionMethod::IS_PUBLIC) as $method) {
            if (Method::match($method)) {
                $methods[$method->getName()] = Method::parse($method);
            }
        }

        return $methods;
    }
}
