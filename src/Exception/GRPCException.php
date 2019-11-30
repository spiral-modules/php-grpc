<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC\Exception;

use Google\Protobuf\Internal\Message;
use Spiral\GRPC\StatusCode;
use Throwable;

class GRPCException extends \RuntimeException
{
    protected const CODE = StatusCode::UNKNOWN;

    /**
     * Collection of protobuf messages for describing error which will be converted to google.protobuf.Any during
     * sending as response.
     *
     * @see https://cloud.google.com/apis/design/errors
     *
     * @var array|Message[]
     */
    private $details;

    /**
     * @param string         $message
     * @param int            $code
     * @param array          $details
     * @param Throwable|null $previous
     */
    public function __construct(string $message = '', int $code = 0, $details = [], Throwable $previous = null)
    {
        if ($code == 0) {
            $code = static::CODE;
        }

        parent::__construct($message, $code, $previous);

        $this->details = $details;
    }

    /**
     * @return array
     */
    public function getDetails(): array
    {
        return $this->details;
    }

    /**
     * @param array $details
     */
    public function setDetails(array $details): void
    {
        $this->details = $details;
    }

    /**
     * Push details message to the exception.
     *
     * @param $details
     *
     * @return $this
     */
    public function withDetails($details)
    {
        $this->details[] = $details;

        return $this;
    }
}
