<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC\Exception;

use Spiral\GRPC\StatusCode;
use Throwable;

class GRPCException extends \RuntimeException
{
    protected const CODE = StatusCode::UNKNOWN;

    /**
     * @param string $message
     * @param int $code
     * @param Throwable|null $previous
     */
    public function __construct(string $message = "", int $code = 0, Throwable $previous = null)
    {
        if ($code == 0) {
            $code = static::CODE;
        }

        parent::__construct($message, $code, $previous);
    }
}