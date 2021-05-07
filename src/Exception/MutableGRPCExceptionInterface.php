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

interface MutableGRPCExceptionInterface extends GRPCExceptionInterface
{
    /**
     * Rewrites details message in the GRPC Exception.
     *
     * @param array<Message> $details
     */
    public function setDetails(array $details): void;

    /**
     * Appends details message to the GRPC Exception.
     *
     * @param Message $message
     */
    public function addDetails(Message $message): void;
}
