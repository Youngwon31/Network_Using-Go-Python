// This code is for the server side.
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// Global logger instance
var log = logrus.New()

type User struct {
	Username string
	Password string
}

var validUser = User{
	Username: "user1",
	Password: "pass1",
}

func handleConnection(conn net.Conn) {
	// Schedule the network connection to be closed via the net.Conn interface.
	defer conn.Close()

	// Using Logrus to log the message
	log.WithFields(logrus.Fields{
		"client": conn.RemoteAddr(),
	}).Info("Client connected") // [2]

	// read the data from client
	reader := bufio.NewReader(conn) // Create a Reader object to read data from the client.

	input, err := reader.ReadString('\n')
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":  err.Error(),
			"client": conn.RemoteAddr(),
		}).Error("Failed to read from client")
		return
	}

	// Log the data received from the client.
	log.WithFields(logrus.Fields{
		"client": conn.RemoteAddr(),
		"data":   input,
	}).Info("Received data from client")

	// Test
	fmt.Printf("Received from client: %s", input)

}

func main() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Set up a log file and include log...
	// study:
	// The os.OpenFile function is the standard Go function for opening files.
	// os.O_CREATE: This flag creates the file if it doesn't exist.
	// os.O_WRONLY:  This flag is for opening files in write-only mode

	// Note: https://github.com/sirupsen/logrus
	logFile, err := os.OpenFile("yseo8834_jiukim9115_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Intended to set log output to logFile.
	log.SetOutput(logFile)

	// UDP > net.Listenpacket, TCP > net.listen
	listener, err := net.Listen("tcp", ":8080") // If an error occurs in this action, it is stored in the err variable.

	fmt.Println("Hello, World!")

	// nil means 'nothing' or 'no error' in Go
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server Started. Listening on port 8080") // [1]

	for {
		fmt.Println("Check [1]")
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Check [2]")
			log.Println("Connection acceptance errors:", err)
			continue
		}

		fmt.Println("Check [3]")
		go handleConnection(conn)
	}

}
