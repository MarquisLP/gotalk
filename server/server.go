package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

// TODO: Ensure messages are displayed in chronological order via synchronization.
// TODO: Handle case where client disconnects using Ctrl+C.

func main() {
	var clientConnections []net.Conn = make([]net.Conn, 0)

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for clients...")

	for {
		connection, err := server.Accept()

		if err != nil {
			fmt.Println("Error while accepting client connection: ", err.Error())
			os.Exit(1)
		}

		clientConnections = append(clientConnections, connection)

		fmt.Println("client connected: " + connection.LocalAddr().String() + ", " + connection.RemoteAddr().String())
		go processClient(connection, &clientConnections)
	}
}

func processClient(connection net.Conn, connections *[]net.Conn) {
	errorOccurred, clientUsername := readMessage(connection)
	if errorOccurred {
		fmt.Println("Failed to receive client username.")
		connection.Close()
		return
	}
	broadcastMessage(clientUsername+" logged in.", connections)

	defer connection.Close()
	for {
		errorOccurred, clientMessage := readMessage(connection)

		if errorOccurred {
			fmt.Println("Error occurred reading message from " + clientUsername + ". Closing connection.")
			break
		}
		if clientMessage == "exit" {
			broadcastMessage(clientUsername+" logged out.", connections)
			break
		}

		go broadcastMessage("["+clientUsername+"] "+clientMessage, connections)
	}
}

func readMessage(connection net.Conn) (bool, string) {
	lengthPrefix := make([]byte, 4)
	_, prefixError := io.ReadFull(connection, lengthPrefix)
	if prefixError == io.EOF {
		return true, ""
	}
	length := binary.BigEndian.Uint32((lengthPrefix))

	messageBytes := make([]byte, int(length))
	_, messageError := io.ReadFull(connection, messageBytes)
	if messageError == io.EOF {
		return true, ""
	}

	message := string(messageBytes[:])

	return false, message
}

func broadcastMessage(message string, connections *[]net.Conn) {
	fmt.Println(message)

	connectionsLength := len(*connections)
	for i := 0; i < connectionsLength; i++ {
		currentConnection := (*connections)[i]
		sendMessage(currentConnection, message)
	}
}

func sendMessage(connection net.Conn, message string) {
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(message)))
	connection.Write(prefix)
	connection.Write([]byte(message))
}
