<?php

/**
 * This file is part of RoadRunner GRPC package.
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

declare(strict_types=1);

namespace Spiral\GRPC\Internal;

/**
 * @internal Json is an internal library class, please do not use it in your code.
 * @psalm-internal Spiral\GRPC
 */
final class Json
{
    /**
     * @var positive-int
     */
    public const DEFAULT_JSON_DEPTH = 512;

    /**
     * @return positive-int|0
     */
    private static function getFlags(): int
    {
        return \defined('\\JSON_THROW_ON_ERROR')
            // Avoid PHP DCE constant inlining
            ? \constant('\\JSON_THROW_ON_ERROR')
            : 0;
    }

    /**
     * @param mixed $payload
     * @return string
     * @throws \JsonException
     */
    public static function encode($payload): string
    {
        /** @var string $result */
        $result = @\json_encode($payload, self::getFlags(), self::DEFAULT_JSON_DEPTH);

        self::assertJsonErrors();

        return $result;
    }

    /**
     * @param string $payload
     * @return array
     * @throws \JsonException
     */
    public static function decode(string $payload): array
    {
        /** @var array $result */
        $result = @\json_decode($payload, true, self::DEFAULT_JSON_DEPTH, self::getFlags());

        self::assertJsonErrors();

        return $result;
    }

    /**
     * @throws \JsonException
     */
    private static function assertJsonErrors(): void
    {
        $code = \json_last_error();

        if ($code !== \JSON_ERROR_NONE) {
            throw new \JsonException(\json_last_error_msg(), $code);
        }
    }
}
