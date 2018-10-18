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
final class Service
{
    /** @var string */
    private $name;

    /** @var object */
    private $handler;

    /** @var InvocatorInterface */
    private $invocator;

    /** @var Method[] */
    private $methods;

    /**
     * @param InvocatorInterface $invocator
     * @param string $interface Service interface name.
     * @param object $handler
     *
     * @throws ServiceException
     */
    public function __construct(InvocatorInterface $invocator, string $interface, $handler)
    {
        $this->invocator = $invocator;
        $this->configure($interface, $handler);
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @return object
     */
    public function getHandler()
    {
        return $this->handler;
    }

    /**
     * @param string $method
     * @param ContextInterface $context
     * @param string $input
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
     * @param string $interface
     * @param object $handler
     *
     * @throws ServiceException
     */
    protected function configure(string $interface, $handler)
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

        if (is_object($handler)) {
            throw new ServiceException("Service handler must be object.");
        }

        if (!$handler instanceof $interface) {
            throw new ServiceException("Service handler does not implement `{$interface}`.");
        }

        $this->handler = $handler;

        // list of all available methods and their object types
        $this->methods = $this->fetchMethods($handler);
    }

    /**
     * @param object $handler
     * @return array
     */
    protected function fetchMethods($handler): array
    {
        $reflection = new \ReflectionObject($handler);

        $methods = [];
        foreach ($reflection->getMethods(\ReflectionMethod::IS_PUBLIC) as $method) {
            if (Method::match($method)) {
                $methods[$method->getName()] = Method::parse($method);
            }
        }

        return $methods;
    }
}