<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;

use Google\Protobuf\Internal\Message;

class Method
{
    private $handler;
    private $name;
    private $input;
    private $output;

    public function __construct(
        object $handler,
        string $name,
        string $input,
        string $output
    ) {
        $this->handler = $handler;
        $this->name = $name;
        $this->input = $input;
        $this->output = $output;
    }

    public function invoke(Context $ctx, string $input): string
    {
        /** @var Message $in */
        $in = new ($this->input);

        try {
            $in->mergeFromString($input);
        } catch (\Exception $e) {
            throw new $e;
        }

        /** @var Message $out */
        $out = call_user_func([$this->handler, $this->name], $ctx, $in);

        try {
            return $out->serializeToString();
        } catch (\Exception $e) {
            throw new $e;
        }
    }
}