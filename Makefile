build:
	go mod tidy
	go fmt
	go test
	go build -o ./bedlamb
