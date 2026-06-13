# kit — batteries-included integrations for harmoni

[![CI](https://github.com/harmonikit/kit/actions/workflows/ci.yml/badge.svg)](https://github.com/harmonikit/kit/actions/workflows/ci.yml)

Multi-module Go repository providing production-ready implementations of harmoni
interfaces. Each directory is an independent Go module — pick what you need, pull
only its dependencies.

## Modules

### Transport

| Module | Import | Dependencies |
|---|---|---|
| HTTP | `github.com/harmonikit/kit/transport/http` | stdlib `net/http` |
| gRPC | `github.com/harmonikit/kit/transport/grpc` | `google.golang.org/grpc` |
| NATS | `github.com/harmonikit/kit/transport/nats` | `github.com/nats-io/nats.go` |
| AMQP | `github.com/harmonikit/kit/transport/amqp` | `github.com/rabbitmq/amqp091-go` |

### Observability

| Module | Import | Dependencies |
|---|---|---|
| Zap | `github.com/harmonikit/kit/log/zap` | `go.uber.org/zap` |
| Zerolog | `github.com/harmonikit/kit/log/zerolog` | `github.com/rs/zerolog` |
| Logrus | `github.com/harmonikit/kit/log/logrus` | `github.com/sirupsen/logrus` |
| Prometheus | `github.com/harmonikit/kit/metrics/prometheus` | `github.com/prometheus/client_golang` |
| OpenTelemetry | `github.com/harmonikit/kit/metrics/otel` | OpenTelemetry SDK |
| OpenTelemetry Tracing | `github.com/harmonikit/kit/tracing/opentelemetry` | `go.opentelemetry.io/otel` |
| Zipkin | `github.com/harmonikit/kit/tracing/zipkin` | Zipkin reporter |

### Infrastructure

| Module | Import | Dependencies |
|---|---|---|
| Consul SD | `github.com/harmonikit/kit/sd/consul` | `github.com/hashicorp/consul/api` |
| Etcd SD | `github.com/harmonikit/kit/sd/etcd` | `go.etcd.io/etcd/client/v3` |
| DNS SD | `github.com/harmonikit/kit/sd/dns` | stdlib `net` |
| Hystrix CB | `github.com/harmonikit/kit/circuitbreaker/hystrix` | — |
| SRE CB | `github.com/harmonikit/kit/circuitbreaker/sre` | — |
| Token Bucket | `github.com/harmonikit/kit/ratelimit/tokenbucket` | — |
| JWT Auth | `github.com/harmonikit/kit/auth/jwt` | `github.com/golang-jwt/jwt/v5` |

### Codec

| Module | Import | Dependencies |
|---|---|---|
| JSON | `github.com/harmonikit/kit/codec/json` | stdlib `encoding/json` |
| Protobuf | `github.com/harmonikit/kit/codec/protobuf` | `google.golang.org/protobuf` |
| MessagePack | `github.com/harmonikit/kit/codec/msgpack` | — |

## Requirements

- Go 1.23+
- `github.com/harmonikit/harmoni` for core interfaces

## License

MIT
