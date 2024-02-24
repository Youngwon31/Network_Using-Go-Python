// This code is for the server side.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Global logger instance
var log = logrus.New()

type User struct {
	Username string
	Password string
}

var validUser = User{
	Username: "jkim9115",
	Password: "9115",
}

func handleConnection(conn net.Conn) {
	// Schedule the network connection to be closed via the net.Conn interface.
	defer conn.Close()

	// read the data from client
	reader := bufio.NewReader(conn) // Create a Reader object to read data from the client.

	for {
		// Using Logrus to log the message
		input, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.WithFields(logrus.Fields{
					"client": conn.RemoteAddr(),
				}).Info("Client connected") // [2]
			} else {
				log.WithFields(logrus.Fields{
					"error":  err.Error(),
					"client": conn.RemoteAddr(),
				}).Error("Failed to read from client")
			}
			return
		}

		// Processing data received from clients  "jkim9115:9115"
		input = strings.TrimSpace(input)
		parts := strings.Split(input, ":")
		if len(parts) != 2 {
			log.WithFields(logrus.Fields{
				"client": conn.RemoteAddr(),
			}).Error("Invalid login format")
			fmt.Fprintln(conn, "Invalid login format")
			continue
		}

		//check
		username := parts[0]
		password := parts[1]

		// Save data received from the client to the log
		log.WithFields(logrus.Fields{
			"client":   conn.RemoteAddr(),
			"username": username,
			"password": password,
		}).Info("Parsed login attempt")

		fmt.Printf("Check\n")
		fmt.Printf("Received username from client: %s\n", username)
		fmt.Printf("Received password client: %s\n", password)

		if username != validUser.Username {
			log.WithFields(logrus.Fields{
				"client": conn.RemoteAddr(),
			}).Warning("Invalid username")
			// fmt.Fprintln(conn, "Invalid username")
		} else if password != validUser.Password {
			log.WithFields(logrus.Fields{
				"client": conn.RemoteAddr(),
			}).Warning("Invalid password")
			// fmt.Fprintln(conn, "Invalid password")
		} else {
			log.WithFields(logrus.Fields{
				"client": conn.RemoteAddr(),
			}).Info("User authenticated successfully")
			// fmt.Fprintln(conn, "Login successful!")
		}

		// Test
		fmt.Printf("Received from client: %s\n", input)
	}

}

func main() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Set up a log file and include log
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
