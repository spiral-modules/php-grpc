<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Google\Protobuf\Internal\Message;
use Spiral\GRPC\Exception\GRPCExceptionInterface;
use Spiral\GRPC\Exception\InvokeException;

final class Invoker implements InvokerInterface
{
    /**
     * @var string
     */
    private const ERROR_METHOD_RETURN =
        'Method %s must return an object that instance of %s, ' .
        'but the result provides type of %s';

    /**
     * @var string
     */
    private const ERROR_METHOD_IN_TYPE =
        'Method %s input type must be an instance of %s, ' .
        'but the input is type of %s'
    ;

    /**
     * {@inheritDoc}
     */
    public function invoke(ServiceInterface $service, Method $method, ContextInterface $ctx, ?string $input): string
    {
        /** @var callable $callable */
        $callable = [$service, $method->getName()];

        /** @var Message $message */
        $message = $callable($ctx, $this->makeInput($method, $input));

        // Note: This validation will only work if the
        // assertions option ("zend.assertions") is enabled.
        assert($this->assertResultType($method, $message));

        try {
            return $message->serializeToString();
        } catch (\Throwable $e) {
            throw InvokeException::create($e->getMessage(), StatusCode::INTERNAL, $e);
        }
    }

    /**
     * Checks that the result from the GRPC service method returns the
     * Message object.
     *
     * @param Method $method
     * @param mixed $result
     * @return bool
     * @throws \BadFunctionCallException
     */
    private function assertResultType(Method $method, $result): bool
    {
        if (! $result instanceof Message) {
            $type = \is_object($result) ? \get_class($result) : \get_debug_type($result);

            throw new \BadFunctionCallException(
                \sprintf(self::ERROR_METHOD_RETURN, $method->getName(), Message::class, $type)
            );
        }

        return true;
    }

    /**
     * @param Method $method
     * @param string|null $body
     * @return Message
     * @throws InvokeException
     */
    private function makeInput(Method $method, ?string $body): Message
    {
        try {
            $class = $method->getInputType();

            // Note: This validation will only work if the
            // assertions option ("zend.assertions") is enabled.
            assert($this->assertInputType($method, $class));

            /** @psalm-suppress UnsafeInstantiation */
            $in = new $class();

            if ($body !== null) {
                $in->mergeFromString($body);
            }

            return $in;
        } catch (\Throwable $e) {
            throw InvokeException::create($e->getMessage(), StatusCode::INTERNAL, $e);
        }
    }

    /**
     * Checks that the input of the GRPC service method contains the
     * Message object.
     *
     * @param Method $method
     * @param string $class
     * @return bool
     * @throws \InvalidArgumentException
     */
    private function assertInputType(Method $method, string $class): bool
    {
        if (! \is_subclass_of($class, Message::class)) {
            throw new \InvalidArgumentException(
                \sprintf(self::ERROR_METHOD_IN_TYPE, $method->getName(), Message::class, $class)
            );
        }

        return true;
    }
}
