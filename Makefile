build:
	templ generate
	go build -o ./bin/webapp ./cmd/webapp/main.go

