CHANGELOG
=========

v1.4.1 (13.10.2020)
-------------------
- RoadRunner version update to 1.8.3
- Golang version in go.mod bump to 1.15
- Add server configuration options (debug) (@aldump)

v1.4.0 (1.08.2020)
-------------------
- Add all major gRPC configuration options. [Docs](https://github.com/spiral/docs/blob/master/grpc/configuration.md#application-server)

v1.3.1 (20.07.2020)
-------------------
- RoadRunner version updated to 1.8.2

v1.3.0 (25.05.2020)
-------------------
- Add the ability to append the certificate authority in the config under the `tls: rootCA` key
- RoadRunner version updated to 1.8.1

v1.2.2 (05.05.2020)
-------------------
- RoadRunner version updated to 1.8.0

v1.2.1 (22.04.2020)
-------------------
- Replaced deprecated github.com/golang/protobuf/proto with new google.golang.org/protobuf/proto
- RoadRunner version updated to 1.7.1

v1.2.0 (27.01.2020)
-------------------
- Add the ability to work on Golang level only (without roadrunner worker and proto file)

v1.1.1 (27.01.2020)
-------------------
- [bugfix] invalid constructor parameters in ServiceException by @everflux

v1.1.0 (30.11.2019)
-------------------
- Add automatic CS fixing
- The minimum PHP version set to 7.2
- Add ResponseHeaders and metadata generation by server by @wolfgang-braun

v1.0.8 (06.09.2019)
-------------------
- Include `limit` and `metrics` service
- Ability to expose GRPC stats to Prometheus

v1.0.7 (22.05.2019)
-------------------
- Server and Invoker are final
- Add support for pool controller (roadrunner 1.4.0) 
- Add strict_types=1

v1.0.4-1.0.6 (26.04.2019)
-------------------
- bugfix, support for imported services in proto annotation by @oneslash 

v1.0.2 (18.03.2019)
-------------------
- Add support for `php_namespace` option
- Add support for nested namespace resolution in generated code
  (thanks to @zarianec)
- protobuf version bump to 1.3.1

v1.0.1 (30.01.2019)
-------------------
- Fix bug causing server not working with empty payloads
- Fix bug with disabled RPC service
- Add elapsed time to the debug log

v1.0.0 (20.10.2018)
-------------------
- initial application release
