build:
	go build ./cmd/consumer

wire:
	wire  ./internal/config/injector

bdd: 
	go test -v ./test/bdd/steps/step_definitions_test.go