<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Spiral\GRPC\Exception\InvokeException;

/**
 * Responsible for data marshalling/unmarshalling and method invocation.
 */
interface InvokerInterface
{
    /**
     * @param ServiceInterface $service
     * @param Method           $method
     * @param ContextInterface $context
     * @param string           $input
     * @return string
     *
     * @throws InvokeException
     */
    public function invoke(
        ServiceInterface $service,
        Method $method,
        ContextInterface $context,
        string $input
    ): string;
}
