<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
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
     * @param Method $method
     * @param ContextInterface $ctx
     * @param string|null $input
     * @return string
     * @throws InvokeException
     */
    public function invoke(ServiceInterface $service, Method $method, ContextInterface $ctx, ?string $input): string;
}
