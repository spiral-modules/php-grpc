<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

/**
 * @template-implements \IteratorAggregate<string, string>
 */
final class ResponseHeaders implements \IteratorAggregate, \Countable
{
    /**
     * @var array<string, string>
     */
    private $headers = [];

    /**
     * @param iterable<string, string> $headers
     */
    public function __construct(iterable $headers = [])
    {
        foreach ($headers as $key => $value) {
            $this->set($key, $value);
        }
    }

    /**
     * @param string $key
     * @param string $value
     */
    public function set(string $key, string $value): void
    {
        $this->headers[$key] = $value;
    }

    /**
     * @param string $key
     * @param string|null $default
     * @return string|null
     */
    public function get(string $key, string $default = null): ?string
    {
        return $this->headers[$key] ?? $default;
    }

    /**
     * {@inheritDoc}
     */
    public function getIterator(): \Traversable
    {
        return new \ArrayIterator($this->headers);
    }

    /**
     * @return int
     */
    public function count(): int
    {
        return \count($this->headers);
    }

    /**
     * @return string
     * @throws \JsonException
     */
    public function packHeaders(): string
    {
        // If an empty array is serialized, it is cast to the string "[]"
        // instead of object string "{}"
        if ($this->headers === []) {
            return '{}';
        }

        $flags = \defined('\\JSON_THROW_ON_ERROR')
            // Avoid PHP DCE constant inlining
            ? \constant('\\JSON_THROW_ON_ERROR')
            : 0;

        return $this->toJsonString($this->headers, $flags);
    }

    /**
     * @param array $payload
     * @param int $flags
     * @return string
     * @throws \JsonException
     */
    private function toJsonString(array $payload, int $flags = 0): string
    {
        $result = @\json_encode($payload, $flags);

        if (($code = \json_last_error()) !== \JSON_ERROR_NONE) {
            throw new \JsonException(\json_last_error_msg(), $code);
        }

        return $result;
    }
}
