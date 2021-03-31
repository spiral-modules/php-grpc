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
use JetBrains\PhpStorm\Deprecated;
use JetBrains\PhpStorm\ExpectedValues;
use Spiral\GRPC\StatusCode;

/**
 * @psalm-import-type StatusCodeType from StatusCode
 */
class GRPCException extends \RuntimeException implements MutableGRPCExceptionInterface
{
    /**
     * Can be overridden by child classes.
     *
     * @psalm-var StatusCodeType
     * @var int
     */
    protected const CODE = StatusCode::UNKNOWN;

    /**
     * Collection of protobuf messages for describing error which will be
     * converted to google.protobuf. Any during sending as response.
     *
     * @see https://cloud.google.com/apis/design/errors
     *
     * @var array<Message>
     */
    private $details;

    /**
     * @param string $message
     * @param StatusCodeType|null $code
     * @param array<Message> $details
     * @param \Throwable|null $previous
     */
    final public function __construct(
        string $message = '',
        #[ExpectedValues(valuesFromClass: StatusCode::class)]
        int $code = null,
        array $details = [],
        \Throwable $previous = null
    ) {
        parent::__construct($message, (int)($code ?? static::CODE), $previous);

        $this->details = $details;
    }

    /**
     * @param string $message
     * @param StatusCodeType|null $code
     * @param array<Message> $details
     * @param \Throwable|null $previous
     * @return static
     */
    public static function create(
        string $message,
        #[ExpectedValues(valuesFromClass: StatusCode::class)]
        int $code = null,
        \Throwable $previous = null,
        array $details = []
    ): self {
        return new static($message, $code, $details, $previous);
    }

    /**
     * {@inheritDoc}
     */
    public function getDetails(): array
    {
        return $this->details;
    }

    /**
     * {@inheritDoc}
     */
    public function setDetails(array $details): void
    {
        $this->details = $details;
    }

    /**
     * {@inheritDoc}
     */
    public function addDetails(Message $message): void
    {
        $this->details[] = $message;
    }

    /**
     * Push details message to the exception.
     *
     * @param Message $details
     * @return $this
     * @deprecated Please use {@see GRPCException::addDetails()} method instead.
     */
    #[Deprecated('Please use GRPCException::addDetails() instead', '%class%::addDetails(%parameter0%)')]
    public function withDetails(
        $details
    ): self {
        $this->details[] = $details;

        return $this;
    }
}
