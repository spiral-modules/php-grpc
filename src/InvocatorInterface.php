<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\GRPCException;
use Spiral\GRPC\Exception\InvokeException;

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
     * @throws InvokeException
     */
    public function invoke($handler, Method $method, ContextInterface $context, string $input): string;
}