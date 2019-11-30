<?php

/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

declare(strict_types=1);

namespace Spiral\GRPC;

final class ResponseHeaders implements \IteratorAggregate
{
    /** @var array */
    private $headers;

    /**
     * @param array $headers
     */
    public function __construct(array $headers)
    {
        $this->headers = $headers;
    }

    /**
     * @param string $key
     * @param mixed  $value
     */
    public function set(string $key, $value): void
    {
        $this->headers[$key] = $value;
    }

    /**
     * @param string $key
     * @param null   $default
     * @return mixed|null
     */
    public function get(string $key, $default = null)
    {
        return $this->headers[$key] ?? $default;
    }

    /**
     * @return array
     */
    public function getIterator(): array
    {
        return $this->headers;
    }

    /**
     * @return string
     */
    public function packHeaders(): string
    {
        return json_encode($this->headers);
    }
}
