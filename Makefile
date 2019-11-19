GOFILES= $$(go list -f '{{join .GoFiles " "}}')

run:
	go run $(GOFILES) -root "https://www.google.com/"

build:
	go build -o $(GOPATH)/bin/crawler $(GOFILES)

run stats:
	go run $(GOFILES) -root "https://www.google.com/" -stats true
