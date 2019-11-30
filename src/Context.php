<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC;

final class Context implements ContextInterface
{
    /** @var array */
    private $values;

    /**
     * @param array $values
     */
    public function __construct(array $values)
    {
        $this->values = $values;
    }

    /**
     * @inheritdoc
     */
    public function withValue(string $key, $value): ContextInterface
    {
        $ctx = clone $this;
        $ctx->values[$key] = $value;

        return $ctx;
    }

    /**
     * @inheritdoc
     */
    public function getValue(string $key)
    {
        return $this->values[$key] ?? null;
    }

    /**
     * @inheritdoc
     */
    public function getValues(): array
    {
        return $this->values;
    }
}
