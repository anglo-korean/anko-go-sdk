// Package anko provides a golang SDK to the Anko Investor Forecasts gRPC service.
//
// It handles authentication and reconnection logic.
//
//	client, err := anko.Connect("anko-1234", "hostname")
//	if err != nil {
//		panic(err)
//	}
//
//	for forecast := range client {
//		// do something with forecasts
//	}
//
//	client.Close()
//
// This SDK accepts an Anko Token and a per-connection name (which may aid debugging where a single token is used across auto-scaled services) and returns a channel of forecasts for ready consumption.
package anko
