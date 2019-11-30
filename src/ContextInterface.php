<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC;

/**
 * Carries information about call context, client information and metadata.
 */
interface ContextInterface
{
    /**
     * Create context with new value.
     *
     * @param string $key
     * @param mixed  $value
     * @return ContextInterface
     */
    public function withValue(string $key, $value): ContextInterface;

    /**
     * Get context value or return null.
     *
     * @param string $key
     * @return mixed|null
     */
    public function getValue(string $key);

    /**
     * Return all context values.
     *
     * @return array
     */
    public function getValues(): array;
}
