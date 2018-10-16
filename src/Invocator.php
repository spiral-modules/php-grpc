<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;


class Invocator
{

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