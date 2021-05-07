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
use Spiral\GRPC\Exception\GRPCException;

/**
 * Method carry information about one specific RPC method, it's input and
 * return types. Provides ability to detect GRPC methods based on given
 * class declaration.
 */
final class Method
{
    /**
     * @var string
     */
    private const ERROR_PARAMS_COUNT =
        'The GRPC method %s can only contain 2 parameters (input and output), but ' .
        'signature contains an %d parameters';

    /**
     * @var string
     */
    private const ERROR_PARAM_UNION_TYPE =
        'Parameter $%s of the GRPC method %s cannot be declared using union type';

    /**
     * @var string
     */
    private const ERROR_PARAM_CONTEXT_TYPE =
        'The first parameter $%s of the GRPC method %s can only take an instance of %s';

    /**
     * @var string
     */
    private const ERROR_PARAM_INPUT_TYPE =
        'The second (input) parameter $%s of the GRPC method %s can only take ' .
        'an instance of %s, but type %s is indicated';

    /**
     * @var string
     */
    private const ERROR_RETURN_UNION_TYPE =
        'Return type of the GRPC method %s cannot be declared using union type';

    /**
     * @var string
     */
    private const ERROR_RETURN_TYPE =
        'Return type of the GRPC method %s must return ' .
        'an instance of %s, but type %s is indicated';

    /**
     * @var string
     */
    private const ERROR_INVALID_GRPC_METHOD = 'Method %s is not valid GRPC method.';

    /**
     * @var string
     */
    private $name;

    /**
     * @var class-string<Message>
     */
    private $input;

    /**
     * @var class-string<Message>
     */
    private $output;

    /**
     * @param string $name
     * @param class-string<Message> $input
     * @param class-string<Message> $output
     */
    private function __construct(string $name, string $input, string $output)
    {
        $this->name = $name;
        $this->input = $input;
        $this->output = $output;
    }

    /**
     * @return string
     */
    public function getName(): string
    {
        return $this->name;
    }

    /**
     * @return class-string<Message>
     */
    public function getInputType(): string
    {
        return $this->input;
    }

    /**
     * @return class-string<Message>
     */
    public function getOutputType(): string
    {
        return $this->output;
    }

    /**
     * @param \ReflectionType|null $type
     * @return \ReflectionClass|null
     * @throws \ReflectionException
     */
    private static function getReflectionClassByType(?\ReflectionType $type): ?\ReflectionClass
    {
        if ($type instanceof \ReflectionNamedType && ! $type->isBuiltin()) {
            /** @psalm-suppress ArgumentTypeCoercion */
            return new \ReflectionClass($type->getName());
        }

        return null;
    }

    /**
     * Returns true if method signature matches.
     *
     * @param \ReflectionMethod $method
     * @return bool
     */
    public static function match(\ReflectionMethod $method): bool
    {
        try {
            self::assertMethodSignature($method);
        } catch (\Throwable $e) {
            return false;
        }

        return true;
    }

    /**
     * @param \ReflectionMethod $method
     * @param \ReflectionParameter $context
     * @throws \ReflectionException
     */
    private static function assertContextParameter(\ReflectionMethod $method, \ReflectionParameter $context): void
    {
        $type = $context->getType();

        // When the type is not specified, it means that it is declared as
        // a "mixed" type, which is a valid case
        if ($type !== null) {
            if (! $type instanceof \ReflectionNamedType) {
                $message = \sprintf(self::ERROR_PARAM_UNION_TYPE, $context->getName(), $method->getName());
                throw new \DomainException($message, 0x02);
            }

            // If the type is not declared as a generic "mixed" or "object",
            // then it can only be a type that implements ContextInterface.
            if (! \in_array($type->getName(), ['mixed', 'object'], true)) {
                /** @psalm-suppress ArgumentTypeCoercion */
                $isContextImplementedType = ! $type->isBuiltin()
                    && (new \ReflectionClass($type->getName()))
                        ->implementsInterface(ContextInterface::class)
                ;

                // Checking that the signature can accept the context.
                //
                // TODO If the type is any other implementation of the Spiral\GRPC\ContextInterface other than
                //      class Spiral\GRPC\Context, it may cause an error.
                //      It might make sense to check for such cases?
                if (! $isContextImplementedType) {
                    $message = \vsprintf(self::ERROR_PARAM_CONTEXT_TYPE, [
                        $context->getName(),
                        $method->getName(),
                        ContextInterface::class
                    ]);

                    throw new \DomainException($message, 0x03);
                }
            }
        }
    }

    /**
     * @param \ReflectionMethod $method
     * @param \ReflectionParameter $input
     * @throws \ReflectionException
     */
    private static function assertInputParameter(\ReflectionMethod $method, \ReflectionParameter $input): void
    {
        $type = $input->getType();

        // Parameter type cannot be omitted ("mixed")
        if ($type === null) {
            $message = \vsprintf(self::ERROR_PARAM_INPUT_TYPE, [
                $input->getName(),
                $method->getName(),
                Message::class,
                'mixed'
            ]);

            throw new \DomainException($message, 0x04);
        }

        // Parameter type cannot be declared as singular non-named type
        if (! $type instanceof \ReflectionNamedType) {
            $message = \sprintf(self::ERROR_PARAM_UNION_TYPE, $input->getName(), $method->getName());
            throw new \DomainException($message, 0x05);
        }

        /** @psalm-suppress ArgumentTypeCoercion */
        $isProtobufMessageType = ! $type->isBuiltin()
            && (new \ReflectionClass($type->getName()))
                ->isSubclassOf(Message::class)
        ;

        if (! $isProtobufMessageType) {
            $message = \vsprintf(self::ERROR_PARAM_INPUT_TYPE, [
                $input->getName(),
                $method->getName(),
                Message::class,
                $type->getName(),
            ]);
            throw new \DomainException($message, 0x06);
        }
    }

    /**
     * @param \ReflectionMethod $method
     * @throws \ReflectionException
     */
    private static function assertOutputReturnType(\ReflectionMethod $method): void
    {
        $type = $method->getReturnType();

        // Return type cannot be omitted ("mixed")
        if ($type === null) {
            $message = \sprintf(self::ERROR_RETURN_TYPE, $method->getName(), Message::class, 'mixed');
            throw new \DomainException($message, 0x07);
        }

        // Return type cannot be declared as singular non-named type
        if (! $type instanceof \ReflectionNamedType) {
            $message = \sprintf(self::ERROR_RETURN_UNION_TYPE, $method->getName());
            throw new \DomainException($message, 0x08);
        }

        /** @psalm-suppress ArgumentTypeCoercion */
        $isProtobufMessageType = ! $type->isBuiltin()
            && (new \ReflectionClass($type->getName()))
                ->isSubclassOf(Message::class)
        ;

        if (! $isProtobufMessageType) {
            $message = \sprintf(self::ERROR_RETURN_TYPE, $method->getName(), Message::class, $type->getName());
            throw new \DomainException($message, 0x09);
        }
    }

    /**
     * @param \ReflectionMethod $method
     * @throws \ReflectionException
     * @throws \DomainException
     */
    private static function assertMethodSignature(\ReflectionMethod $method): void
    {
        // Check that there are only two parameters
        if ($method->getNumberOfParameters() !== 2) {
            $message = \sprintf(self::ERROR_PARAMS_COUNT, $method->getName(), $method->getNumberOfParameters());
            throw new \DomainException($message, 0x01);
        }

        /**
         * @var \ReflectionParameter $context
         * @var \ReflectionParameter $input
         */
        [$context, $input] = $method->getParameters();

        // The first parameter can only take a context object
        self::assertContextParameter($method, $context);

        // The second argument can only be a subtype of the Google\Protobuf\Internal\Message class
        self::assertInputParameter($method, $input);

        // The return type must be declared as a Google\Protobuf\Internal\Message class
        self::assertOutputReturnType($method);
    }

    /**
     * Creates a new {@see Method} object from a {@see \ReflectionMethod} object.
     *
     * @param \ReflectionMethod $method
     * @return Method
     */
    public static function parse(\ReflectionMethod $method): Method
    {
        try {
            self::assertMethodSignature($method);
        } catch (\Throwable $e) {
            $message = \sprintf(self::ERROR_INVALID_GRPC_METHOD, $method->getName());
            throw GRPCException::create($message, StatusCode::INTERNAL, $e);
        }

        [,$input] = $method->getParameters();

        /** @var \ReflectionNamedType $inputType */
        $inputType = $input->getType();

        /** @var \ReflectionNamedType $returnType */
        $returnType = $method->getReturnType();

        /** @psalm-suppress ArgumentTypeCoercion */
        return new self($method->getName(), $inputType->getName(), $returnType->getName());
    }
}
