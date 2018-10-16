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

/**
 * Wraps handlers methods.
 */
class Service
{
    /** @var string */
    private $name;

    /** @var InvocatorInterface */
    private $invocator;

    /** @var object */
    private $handler;

    /** @var Method[] */
    private $methods;

    /**
     * @param string             $name
     * @param InvocatorInterface $invocator
     * @param object             $handler
     */
    public function __construct(string $name, InvocatorInterface $invocator, $handler)
    {
        $this->name = $name;
        $this->invocator = $invocator;
        $this->handler = $handler;

        // list of all available methods and their object types
        $this->methods = $this->fetchMethods($handler);
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
     * @param string           $method
     * @param ContextInterface $context
     * @param string           $input
     * @return string
     *
     * @throws InvokeException
     */
    public function invoke(string $method, ContextInterface $context, string $input): string
    {
        if (!isset($this->methods[$method])) {
            throw new NotFoundException("Method `{$method}` not found in service `{$this->name}`.");
        }

        return $this->invocator->invoke(
            $this->handler,
            $this->methods[$method],
            $context,
            $input
        );
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