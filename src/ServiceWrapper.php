<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

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
    private $handler;

    /** @var InvokerInterface */
    private $invocator;

    /** @var Method[] */
    private $methods;

    /**
     * @param InvokerInterface $invocator
     * @param string           $interface Service interface name.
     * @param ServiceInterface $service
     *
     * @throws ServiceException
     */
    public function __construct(
        InvokerInterface $invocator,
        string $interface,
        ServiceInterface $service
    ) {
        $this->invocator = $invocator;
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
        return $this->handler;
    }

    /**
     * @return array
     */
    public function getMethods(): array
    {
        return array_values($this->methods);
    }

    /**
     * @param string           $method
     * @param ContextInterface $context
     * @param string           $input
     * @return string
     *
     * @throws NotFoundException
     * @throws InvokeException
     */
    public function invoke(string $method, ContextInterface $context, string $input): string
    {
        if (!isset($this->methods[$method])) {
            throw new NotFoundException("Method `{$method}` not found in service `{$this->name}`.");
        }

        return $this->invocator->invoke($this->handler, $this->methods[$method], $context, $input);
    }

    /**
     * Configure service name and methods.
     *
     * @param string           $interface
     * @param ServiceInterface $service
     *
     * @throws ServiceException
     */
    protected function configure(string $interface, ServiceInterface $service)
    {
        try {
            $r = new \ReflectionClass($interface);
            if (!$r->hasConstant('NAME')) {
                throw new ServiceException(
                    "Invalid service interface `{$interface}`, constant `NAME` not found."
                );
            }
            $this->name = $r->getConstant('NAME');
        } catch (\ReflectionException $e) {
            throw new ServiceException(
                "Invalid service interface `{$interface}`.",
                StatusCode::INTERNAL,
                $e
            );
        }

        if (!is_object($service)) {
            throw new ServiceException("Service handler must be an object.");
        }

        if (!$service instanceof $interface) {
            throw new ServiceException("Service handler does not implement `{$interface}`.");
        }

        $this->handler = $service;

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