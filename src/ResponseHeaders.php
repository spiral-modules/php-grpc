<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC;

use Spiral\GRPC\Internal\Json;

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

        return Json::encode($this->headers);
    }
}
