<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;


interface InvocatorInterface
{
    public function invoke($handler, Method $method, ContextInterface $context, string $input): string;
}