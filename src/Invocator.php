<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */

namespace Spiral\GRPC;


use Google\Protobuf\Internal\Message;

class Invocator implements InvocatorInterface
{
    // todo: map exceptions
    public function invoke($handler, Method $method, ContextInterface $context, string $input): string
    {
        try {
            $in = $this->makeInput($method);
            $in->mergeFromString($input);
        } catch (\Throwable $e) {
            throw new $e;
        }

        try {
            $out = call_user_func([$handler, $method->getName()], $context, $in);
        } catch (\Throwable $e) {
            throw new $e;
        }

        try {
            return $out->serializeToString();
        } catch (\Throwable $e) {
            throw new $e;
        }
    }

    /**
     * @param Method $method
     * @return Message
     */
    private function makeInput(Method $method): Message
    {
        $in = $method->getInputType();

        return new $in;
    }
}