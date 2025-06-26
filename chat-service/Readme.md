# Go Chat Service

This project is a Go-based chat service. Below are the instructions to set up, build, and run the project using the provided Makefile.

## Prerequisites

- Go (https://golang.org/dl/)
- Make (https://www.gnu.org/software/make/)

## Setup

To set up the project, run the following command:

```sh
make setup
```

## Build

To build the project, run:

```sh
make build
```

The binary will be created in the `./dist` directory.

## Running the Server

### HTTP Server

To start the HTTP server, run:

```sh
make http-serve
```

The server will start on port 8081.

### WebSocket Server

To start the WebSocket server, run:

```sh
make ws-serve
```

The server will start on the default WebSocket port.

## Generating Mocks

To generate mocks using `mockery`, run:

```sh
make generate-mock
```
