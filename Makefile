DOCUMENTATION_CHECKOUT ?= ../documentation

gateway.pb.go gateway_grpc.pb.go: $(DOCUMENTATION_CHECKOUT)/proto/gateway.proto
	protoc -I $(dir $<) $< --go_out=module=github.com/anglo-korean/anko-go-sdk:. --go-grpc_out=module=github.com/anglo-korean/anko-go-sdk:.

README.md: doc.go
	goreadme -import-path github.com/anglo-korean/anko-go-sdk -title github.com/anglo-korean/anko-go-sdk -badge-godoc -badge-goreportcard -skip-sub-packages > $@
