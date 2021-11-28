package main

import (
	"log"
	"os"

	"github.com/anglo-korean/anko-go-sdk"
)

var Token = os.Getenv("ANKO_TOKEN")

func main() {
	client, err := anko.New(Token, "my-client")
	if err != nil {
		panic(err)
	}

	panic(client.Handle(handler))
}

func handler(f *anko.Forecast) error {
	log.Printf("ID: %s", f.Id)
	log.Printf("Symbol: %v", f.Symbol)
	log.Printf("Label: %s", f.Label.String())
	log.Printf("Confidence Score: %.3f", f.Score)
	log.Print()

	return nil
}
