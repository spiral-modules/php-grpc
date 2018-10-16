<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;


interface ContextInterface
{
    public function withValue(string $key, $value): ContextInterface;

    public function getValue(string $key);
}