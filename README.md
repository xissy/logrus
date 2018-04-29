# logrus

## JSONFormatter for Apex Up

* Original JSONFormatter

```json
{
  "grpc.code":"OK",
  "grpc.method":"LogIn",
  "grpc.service":"pb.taeho.account.AccountService",
  "grpc.start_time":"2018-03-30T18:02:25-07:00",
  "grpc.time_ms":3.615,
  "level":"info",
  "msg":"finished unary call with code OK",
  "peer.address":"127.0.0.1:62832",
  "span.kind":"server",
  "system":"grpc"
}
```

* Apex Up JSONFormatter

```json
{
  "fields":{
    "grpc.code":"OK",
    "grpc.method":"Auth",
    "grpc_service":"pb.taeho.auth.AuthService",
    "grpc_start_time":"2018-03-30T10:55:34-07:00",
    "grpc_time_ms":0.25,
    "peer_address":"127.0.0.1:61161",
    "span_kind":"server",
    "system":"grpc"
  },
  "level":"info",
  "message":"finished unary call with code OK"
}
```

## Usage

```go
package main

import (
    log "github.com/sirupsen/logrus"
    "github.com/xissy/logrus"
)

func main() {
    log.SetOutput(os.Stdout)
    log.SetFormatter(&logrus.ApexUpJSONFormatter{})
    log.WithField("key", "value").Info("message here")
}
```

## References

* https://github.com/sirupsen/logrus
* https://github.com/apex/up
* https://github.com/grpc-ecosystem/go-grpc-middleware
