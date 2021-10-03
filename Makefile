##
## Tools
##
tools:
	@mkdir -p ./bin
	@go get -u github.com/cosmtrek/air
	@go mod tidy

##
## Start Application
##
run: tools
	air

##
## Tests
##
tests:
	go test -cover ./...

##
## Build application
##
build:
	@mkdir -p ./bin
	GO111MODULE=on GOGC=off go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -i -o ./bin/server ./main.go