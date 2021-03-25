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
 * Carries information about call context, client information and metadata.
 */
interface ContextInterface
{
    /**
     * Create context with new value.
     *
     * @param string $key
     * @param mixed $value
     * @return $this
     */
    public function withValue(string $key, $value): self;

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
     * @return array<string, mixed>
     */
    public function getValues(): array;
}
