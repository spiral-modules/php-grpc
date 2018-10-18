<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Google\Protobuf\Internal\Message;

/**
 * Carries information about service method.
 */
final class Method
{
    /** @var string */
    private $name;

    /** @var string */
    private $input;

    /** @var string */
    private $output;

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @return string
     */
    public function getInputType(): string
    {
        return $this->input;
    }

    /**
     * @return string
     */
    public function getOutputType(): string
    {
        return $this->output;
    }

    /**
     * Returns true if method signature matches.
     *
     * @param \ReflectionMethod $method
     * @return bool
     */
    public static function match(\ReflectionMethod $method): bool
    {
        if ($method->getNumberOfParameters() != 2) {
            return false;
        }

        $ctx = $method->getParameters()[0]->getClass();
        $in = $method->getParameters()[1]->getClass();

        if (empty($ctx) || !$ctx->implementsInterface(ContextInterface::class)) {
            return false;
        }

        if (empty($in) || !$in->isSubclassOf(Message::class)) {
            return false;
        }

        if (empty($method->getReturnType())) {
            return false;
        }

        try {
            $return = $method->getReturnType()->getName();
            if (!class_exists($return)) {
                return false;
            }
            $return = new \ReflectionClass($return);
        } catch (\ReflectionException $e) {
            return false;
        }

        return $return->isSubclassOf(Message::class);
    }

    /**
     * Returns true if method signature matches.
     *
     * @param \ReflectionMethod $method
     * @return Method
     */
    public static function parse(\ReflectionMethod $method): Method
    {
        $m = new self;
        $m->name = $method->getName();
        $m->input = $method->getParameters()[1]->getClass()->getName();
        $m->output = $method->getReturnType()->getName();

        return $m;
    }
}