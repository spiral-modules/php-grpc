PHP-GRPC
=================================
[![Latest Stable Version](https://poser.pugx.org/spiral/grpc/version)](https://packagist.org/packages/spiral/grpc)
[![GoDoc](https://godoc.org/github.com/spiral/php-grpc?status.svg)](https://godoc.org/github.com/spiral/php-grpc)
[![Build Status](https://travis-ci.org/spiral/php-grpc.svg?branch=master)](https://travis-ci.org/spiral/php-grpc)
[![Go Report Card](https://goreportcard.com/badge/github.com/spiral/php-grpc)](https://goreportcard.com/report/github.com/spiral/php-grpc)
[![Codecov](https://codecov.io/gh/spiral/php-grpc/branch/master/graph/badge.svg)](https://codecov.io/gh/spiral/php-grpc/)

PHP-GRPC is an open source (MIT licensed) high-performance PHP GRPC server build at top of [RoadRunner](https://github.com/spiral/roadrunner).
Server support both PHP and Golang services running within one application. 

Features:
--------
- comliant with native Golang GRPC services
- very fast, minimal proxy overlay
- easy TLS configuration
- debug tools included
- middleware and server options support
- code generation using `protoc` plugin
- transport, message, worker error management
- response error codes over php exceptions
- works on Windows

License:
--------
The MIT License (MIT). Please see [`LICENSE`](./LICENSE) for more information. Maintained by [SpiralScout](https://spiralscout.com).
