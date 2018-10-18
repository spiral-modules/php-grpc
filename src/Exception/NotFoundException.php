<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC\Exception;

use Spiral\GRPC\StatusCode;

class NotFoundException extends InvokeException
{
    protected const CODE = StatusCode::NOT_FOUND;
}