<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

/**
 * @template-implements \IteratorAggregate<string, mixed>
 * @template-implements \ArrayAccess<string, mixed>
 */
final class Context implements ContextInterface, \IteratorAggregate, \Countable, \ArrayAccess
{
    /**
     * @var array<string, mixed>
     */
    private $values;

    /**
     * @param array<string, mixed> $values
     */
    public function __construct(array $values)
    {
        $this->values = $values;
    }

    /**
     * {@inheritDoc}
     */
    public function withValue(string $key, $value): ContextInterface
    {
        $ctx = clone $this;
        $ctx->values[$key] = $value;

        return $ctx;
    }

    /**
     * {@inheritDoc}
     * @param mixed|null $default
     */
    public function getValue(string $key, $default = null)
    {
        return $this->values[$key] ?? $default;
    }

    /**
     * {@inheritDoc}
     */
    public function getValues(): array
    {
        return $this->values;
    }

    /**
     * {@inheritDoc}
     */
    public function offsetExists($offset): bool
    {
        assert(\is_string($offset), 'Offset argument must be a type of string');

        /**
         * Note: PHP Opcode optimisation
         * @see https://www.php.net/manual/pt_BR/internals2.opcodes.isset-isempty-var.php
         *
         * Priority use `ZEND_ISSET_ISEMPTY_VAR !0` opcode instead of `DO_FCALL 'array_key_exists'`.
         */
        return isset($this->values[$offset]) || \array_key_exists($offset, $this->values);
    }

    /**
     * {@inheritDoc}
     */
    public function offsetGet($offset)
    {
        assert(\is_string($offset), 'Offset argument must be a type of string');

        return $this->values[$offset] ?? null;
    }

    /**
     * {@inheritDoc}
     */
    public function offsetSet($offset, $value): void
    {
        assert(\is_string($offset), 'Offset argument must be a type of string');

        $this->values[$offset] = $value;
    }

    /**
     * {@inheritDoc}
     */
    public function offsetUnset($offset): void
    {
        assert(\is_string($offset), 'Offset argument must be a type of string');

        unset($this->values[$offset]);
    }


    /**
     * {@inheritDoc}
     */
    public function getIterator(): \Traversable
    {
        return new \ArrayIterator($this->values);
    }

    /**
     * {@inheritDoc}
     */
    public function count(): int
    {
        return \count($this->values);
    }
}
