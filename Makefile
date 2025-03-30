build:
	npx @tailwindcss/cli -i ./internal/webapp/ui/assets/css/input.css -o ./internal/webapp/ui/assets/css/output.css
	templ generate
	go build -o ./bin/webapp ./cmd/webapp/main.go

