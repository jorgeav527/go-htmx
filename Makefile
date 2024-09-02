BINARY_NAME=my-app

build:
	./tailwindcss -i views/css/styles.css -o public/styles.css
	@templ generate view
	@go build -o bin/fullstackgo main.go 

test:
	@go test -v ./...
	
run: build
	@./bin/fullstackgo

tailwind:
	@./tailwindcss -i views/css/styles.css -o public/styles.css --watch

templ:
	@templ generate -watch -proxy=http://localhost:8080


clean:
	rm -f ./bin/$(BINARY_NAME)*

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug main.go