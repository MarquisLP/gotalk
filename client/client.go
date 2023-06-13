package main

import (
	"bufio"
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

func main() {
	fmt.Printf("Connecting to %s server on %s:%s...\n", SERVER_TYPE, SERVER_HOST, SERVER_PORT)
	connection := connectToServer(SERVER_TYPE, SERVER_HOST, SERVER_PORT)
	fmt.Println("Connection successful!")

	fmt.Print("Username: ")
	username := getLineFromStdin()

	sendToConnection(connection, username)

	receiveChannel := make(chan bool)
	go sendMessagesUntilExit(username, connection)
	go receiveMessagesUntilExit(connection, receiveChannel)
	select {
	case <-receiveChannel:
		os.Exit(0)
		break
	}
}

func connectToServer(serverType string, host string, port string) net.Conn {
	connection, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		panic(err)
	}

	return connection
}

func getLineFromStdin() string {
	userInputScanner := bufio.NewScanner(os.Stdin)
	userInputScanner.Scan()
	return userInputScanner.Text()
}

func sendToConnection(connection net.Conn, message string) {
	prefix := make([]byte, 4)
	binary.BigEndian.PutUint32(prefix, uint32(len(message)))
	connection.Write(prefix)
	connection.Write([]byte(message))
}

func sendMessagesUntilExit(username string, connection net.Conn) {
	for {
		userMessage := getLineFromStdin()
		sendToConnection(connection, userMessage)
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

func receiveMessagesUntilExit(connection net.Conn, channel chan bool) {
	for {
		errorOccurred, message := readMessage(connection)
		if errorOccurred {
			fmt.Println("Server closed connection.")
			break
		}
		fmt.Println(message)
	}
	channel <- true
}
