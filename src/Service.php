<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\MethodNotFoundException;

class Service
{
    /** @var string */
    private $name;

    /** @var object */
    private $handler;

    /** @var Method[] */
    private $methods = [];

    /**
     * @param string $name
     * @param object $handler
     */
    public function __construct(string $name, object $handler)
    {
        $this->name = $name;
        $this->handler = $handler;

        $this->fetchMethods($handler);
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
    public function getHandler(): object
    {
        return $this->handler;
    }

    /**
     * @return array
     */
    public function getMethods(): array
    {
        return $this->methods;
    }

    /**
     * @param string  $method
     * @param Context $context
     * @param string  $input
     * @return string
     *
     * @throws MethodNotFoundException
     */
    public function invoke(string $method, Context $context, string $input): string
    {
        if (!isset($this->methods[$method])) {
            throw new MethodNotFoundException(
                "Method `{$method}` not found in service `{$this->name}`."
            );
        }

        return $this->methods[$method]->invoke($context, $input);
    }

    /**
     * @param object $handler
     */
    protected function fetchMethods(object $handler)
    {
        $reflection = new \ReflectionObject($handler);

        foreach ($reflection->getMethods(\ReflectionMethod::IS_PUBLIC) as $method) {
            dumP($method);
        }
    }
}