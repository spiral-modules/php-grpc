<?php
/**
 * Spiral Framework.
 *
 * @license   MIT
 * @author    Anton Titov (Wolfy-J)
 */
require 'GPBMetadata/Test.php';
require 'Message.php';

$t = new Message();
$t->setMsg("message")->setValue(190000);

file_put_contents("data.dat", $t->serializeToString());