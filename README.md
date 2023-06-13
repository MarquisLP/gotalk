# gotalk

This is a simple chat server and client implementation written in Go. It leverages goroutines to easily and efficiently broadcast chat messages from one user to all other users.

The server and client both run locally on the same machine and communicate over HTTP. Future developments could include cross-machine communication.

## Setup

1. In a terminal instance, start the server by running: `go run ./server/server.go`
2. For each chat user, open another terminal instance and run: `go run ./server/server.go`
3. To log out as a user, type "exit" (all lowercase)
4. To close the server, use Ctrl+C.
