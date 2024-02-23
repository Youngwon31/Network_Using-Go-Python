// This code is for the server side.

package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// Global logger instance
var log = logrus.New()

func main() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Set up a log file and include log...
	// study:
	// The os.OpenFile function is the standard Go function for opening files.
	// os.O_CREATE: This flag creates the file if it doesn't exist.
	// os.O_WRONLY:  This flag is for opening files in write-only mode

	// Note: https://github.com/sirupsen/logrus
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

	log.Println("Server Started. Listening on port 8080")

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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Using Logrus to log the message
	log.WithFields(logrus.Fields{
		"client": conn.RemoteAddr(),
	}).Info("Client connected")

	// 메시지 읽기 및 로깅 로직을 여기에 구현...
}
