<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */
declare(strict_types=1);

namespace Spiral\GRPC;

final class Context implements ContextInterface
{
    /** @var array */
    private $values;

    /** @var array */
    private $outgoingHeaders;

    /**
     * @param array $values
     * @param array $outgoingHeaders
     */
    public function __construct(array $values, ?array $outgoingHeaders = [])
    {
        $this->values = $values;
        $this->outgoingHeaders = $outgoingHeaders;
    }

    /**
     * @inheritdoc
     */
    public function withValue(string $key, $value): ContextInterface
    {
        $ctx = clone $this;
        $ctx->values[$key] = $value;

        return $ctx;
    }

    /**
     * @inheritdoc
     */
    public function getValue(string $key)
    {
        return $this->values[$key] ?? null;
    }

    /**
     * @inheritdoc
     */
    public function getValues(): array
    {
        return $this->values;
    }

    /**
     * @inheritdoc
     */
    public function getOutgoingHeader(string $key)
    {
        return $this->outgoingHeaders[$key] ?? null;
    }

    /**
     * @inheritdoc
     */
    public function getOutgoingHeaders(): array
    {
        return $this->outgoingHeaders;
    }

    /**
     * @inheritdoc
     */
    public function appendOutgoingHeader(array $headers): void
    {
        $this->outgoingHeaders = array_merge($this->outgoingHeaders, $headers);
    }
}