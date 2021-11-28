// Package anko provides a golang SDK to the Anko Investor Forecasts gRPC service.
//
// It handles authentication and reconnection logic.
//
//	client, err := anko.New("anko-1234", "my-client")
//	if err != nil {
//		panic(err)
//	}
//
//	panic(client.Handle(func(f *anko.Forecast) error {
//		log.Printf("%#v", f)
//		return nil
//	}))
//
//
// This SDK accepts an Anko Token and a per-connection name (which may aid debugging where a single token is used across auto-scaled services) and returns a channel of forecasts for ready consumption.
package anko
