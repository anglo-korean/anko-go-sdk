# github.com/anglo-korean/anko-go-sdk

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/anglo-korean/anko-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/anglo-korean/anko-go-sdk)](https://goreportcard.com/report/github.com/anglo-korean/anko-go-sdk)

Package anko provides a golang SDK to the Anko Investor Forecasts gRPC service.

```go
import "github.com/anglo-korean/anko-go-sdk"
```

It handles authentication and reconnection logic.

```go
anko.ConnectionTimeout = time.Second // default: 5s

client, err := anko.New("anko-1234", "my-client")
if err != nil {
	panic(err)
}

panic(client.Handle(func(f *anko.Forecast) error {
	log.Printf("%#v", f)
	return nil
}))
```

This SDK accepts an Anko Token and a per-connection name (which may aid debugging where a single token is used across auto-scaled services) and returns a channel of forecasts for ready consumption.

client.Handle accepts an anko.Handler function, a convention borrowed from packages such as net/http.
This function looks like:

```go
type Handler func(*anko.Forecast) error
```

This function will be called on every Forecast received from the Anko Investor Forecasts gRPC Service.
Any returned error will halt the processing of Forecasts. By default these functions are called sequentially; this is the opposite behaviour to http.HandlerFunc, where each handler runs in a new go function.

If this behaviour is important to you, then simply write your anko.Handler to provide its own sync/ dispatch/ gofunc logic.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
