<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\GRPCException;

/**
 * Responsible for data marshalling/unmarshalling and method invocation.
 */
interface InvocatorInterface
{
    /**
     * @param object $handler
     * @param Method $method
     * @param ContextInterface $context
     * @param string $input
     * @return string
     *
     * @throws GRPCException
     */
    public function invoke($handler, Method $method, ContextInterface $context, string $input): string;
}