<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC\Exception;

use Google\Protobuf\Internal\Message;
use JetBrains\PhpStorm\ExpectedValues;
use Spiral\GRPC\StatusCode;

/**
 * @psalm-import-type StatusCodeType from StatusCode
 */
interface GRPCExceptionInterface extends \Throwable
{
    /**
     * Returns GRPC exception status code.
     *
     * @psalm-suppress MissingImmutableAnnotation
     * @psalm-return StatusCodeType
     * @return int
     */
    #[ExpectedValues(valuesFromClass: StatusCode::class)]
    public function getCode();

    /**
     * @return array<Message>
     */
    public function getDetails(): array;
}
