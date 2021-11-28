package anko

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const addr = "forecasts.anko-investor.com:443"
const ua = "github.com/anglo-korean/anko-go-sdk#0.1.0"

var errBadConnection = fmt.Errorf("unable to connect to service on time")

// ConnectionTimeout is used in two places:
// * To timeout connections to the Forecasts gRPC Service, and
// * To provide a time limit for the Forecasts gRPC Service to validate an anko token and signal readiness
//
// This timeout allows consumer applications to recognise when a gRPC connection is hanging, and when
// a lack of Forecasts just means there are no valid Forecasts.
var ConnectionTimeout = time.Second * 5

// Handler is a function, in much the same vein as http.HandlerFunc, which consumer
// applications may use to process Forecasts.
//
// Handler is called on every received forecasts. Any error returned by it will halt
// the consumer, and so Handlers must make sure to only return errors where this
// behaviour is appropriate.
//
// Similarly Handlers are called synchonously, which is the opposite behaviour to
// how http.HandlerFunc works- it is the responsibility of the developer to provide
// gofunc/ sync semantics where required.
type Handler func(*Forecast) error

// Connection represents a connection to the Anko Investor Forecasts gRPC service.
//
// Once instantiated, it is used to stream Forecasts for consumer applications to use
// as they want.
type Connection struct {
	client     ForecastsClient
	token      string
	identifier string
}

// New creates a connection to the Anko Investor gRPC service
func New(token, identifier string) (c Connection, err error) {
	c.token = token
	c.identifier = identifier

	err = c.connect()

	return
}

func (c *Connection) connect() (err error) {
	ctx, _ := context.WithTimeout(context.Background(), ConnectionTimeout)

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithBlock())
	if err != nil {
		return
	}

	c.client = NewForecastsClient(conn)

	return
}

// Handle accepts a Handler in order to provide Forecasts to consumer applications.
//
// As detailed in the documentation for Handler, Forecasts are processed one-at-a-time,
// when the previous Forecast has been processed.
//
// This behaviour may be avoided by having your Handler `go func` out to another function.
//
// Any errors returned from a Handler will cause this function to return.
func (c Connection) Handle(handler Handler) (err error) {
	for {
		err = c.handler(handler)
		if err != nil && err.Error() == "rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: INTERNAL_ERROR" {
			err = c.connect()
			if err != nil {
				break
			}

			continue
		}

		break
	}

	return
}

func (c Connection) handler(handler Handler) (err error) {
	ctx := context.Background()
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("bearer %s", c.token)})
	ctx = metadata.NewOutgoingContext(ctx, md)

	m := &Metadata{
		Ua: ua,
		Tags: []*Tag{
			{Key: "Identifier", Value: c.identifier},
		},
	}

	sc, err := c.client.Stream(ctx, m)
	if err != nil {
		return
	}

	err = c.testConn(sc)
	if err != nil {
		return
	}

	var (
		f *Forecast
	)

	for {
		f, err = sc.Recv()
		if err != nil {
			return
		}

		err = handler(f)
		if err != nil {
			return
		}
	}
}

// testConn will stream a single messgae from the Stream endpoint, returning
// an error if either the Recv call fails, or if the wrong RIC is returned
func (c Connection) testConn(sc Forecasts_StreamClient) (err error) {
	out := make(chan error)
	go func() {
		defer close(out)

		f, err := sc.Recv()
		if err != nil {
			out <- err

			return
		}

		if f.Ric != "DUMMY" {
			out <- errBadConnection

			return
		}

		out <- nil
	}()

	// On connect, the Forecasts service will send a dummy forecast to show
	// that the connection is up and running.
	//
	// Wait for this to either turn up, or timeout
	select {
	case err = <-out:
		if err != nil {
			return
		}

	case <-time.After(ConnectionTimeout):
		err = errBadConnection

		return
	}

	return
}
