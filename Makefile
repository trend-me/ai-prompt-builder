build:
	go mod tidy
	go build -o consumer.out ./cmd/consumer

wire:
	wire  ./internal/config/injector

bdd: 
	go test -v ./test/bdd/steps/step_definitions_test.go