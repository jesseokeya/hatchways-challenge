# Hatchways Challenge

Backend service written with love

### Prerequisites

Make sure you have golang installed on your system. [Golang Instalation instruction](https://golang.org/doc/install)


### Installing Dependencies

To run this project you will need to install the third party dependencies.

Lists dependencies

```
go list -m all
```

Install dependencies with

```
go mod download
```

### Runing Application Locally

To run this application use the command below

```
make run
```

### Runing In Docker
To run this application using docker use the following command(s) below

```
docker-compose up -d 
```

### Running Binary File

To run this application use the command below

```
sudo make build
```

And

```
./bin/server
```

## Running the tests

To run the automated tests for this system

```
go test ./...
```

## Built With
* [Chi](https://github.com/go-chi/chi) - The web framework used

## Authors
* **Jesse Okeya**