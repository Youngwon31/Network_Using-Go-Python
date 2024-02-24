/*
* FILE          : Server.go
* PROJECT       : Assignment3 - Services and Logging
* AUTHOR        : Youngwon Seo(8818834), Jiu Kim(8819115)
* DATE          : 2024.02.24
* DESCRIPTION   : This program is a server-side application that uses TCP to establish a connection with a client.
*                 It handles requests, such as login attempts from clients, and logs various events
*                 such as successful authentication, invalid login attempts, rate limit exceeded, etc.
*                 The server is designed to handle multiple concurrent client connections, allowing it to take care of enforcing rate limits to prevent abuse
*                 and accurately track and log each client's activity for review.
 */

// This code is for the server side.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Global logger instance
// Note: https://github.com/sirupsen/logrus
var log = logrus.New()

type User struct {
	Username string
	Password string
}

var validUser = User{
	Username: "jkim9115",
	Password: "9115",
}

// requestTracker was created for the purpose of tracking the time of a client's request.
var requestTracker = make(map[string][]time.Time)

func handleConnection(conn net.Conn, clientID int) {
	// Schedule the network connection to be closed via the net.Conn interface.
	defer conn.Close()

	// read the data from client
	reader := bufio.NewReader(conn) // Create a Reader object to read data from the client.
	clientAddr := conn.RemoteAddr().String()

	for {
		// Using Logrus to log the message
		input, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.WithFields(logrus.Fields{
					"client":   conn.RemoteAddr().String(),
					"clientID": clientID,
				}).Info("Client connected")
			} else {
				log.WithFields(logrus.Fields{
					"error":    err.Error(),
					"client":   conn.RemoteAddr().String(),
					"clientID": clientID,
				}).Error("Failed to read from client")
			}
			return
		}

		// Update request tracking based on the current time
		updateRequestTracker(clientAddr)

		// Check for more than 10 requests in 10 seconds
		if len(requestTracker[clientAddr]) > 10 {
			log.WithFields(logrus.Fields{
				"client":   conn.RemoteAddr().String(),
				"clientID": clientID,
			}).Warning("Rate limit exceeded, disconnecting client")
			return
		}

		// Processing data received from clients  "jkim9115:9115"
		input = strings.TrimSpace(input)
		parts := strings.Split(input, ":")
		if len(parts) != 2 {
			log.WithFields(logrus.Fields{
				"client":   conn.RemoteAddr().String(),
				"clientID": clientID,
			}).Error("Invalid login format")
			continue
		}

		//check
		username := parts[0]
		password := parts[1]

		// Save data received from the client to the log
		log.WithFields(logrus.Fields{
			"client":   conn.RemoteAddr().String(),
			"clientID": clientID,
			"username": username,
			"password": password,
		}).Info("Login attempt")

		fmt.Printf("Received username from client[%d]: %s\n", clientID, username)
		fmt.Printf("Received password from client[%d]: %s\n", clientID, password)

		if username != validUser.Username {
			log.WithFields(logrus.Fields{
				"client":   conn.RemoteAddr().String(),
				"clientID": clientID,
				"username": username,
				"password": password,
			}).Warning("Invalid username")
			// fmt.Fprintln(conn, "Invalid username")
		} else if password != validUser.Password {
			log.WithFields(logrus.Fields{
				"client":   conn.RemoteAddr().String(),
				"clientID": clientID,
				"username": username,
				"password": password,
			}).Warning("Invalid password")
			// fmt.Fprintln(conn, "Invalid password")
		} else {
			log.WithFields(logrus.Fields{
				"client":   conn.RemoteAddr().String(),
				"clientID": clientID,
				"username": username,
				"password": password,
			}).Info("User authenticated successfully")
			// fmt.Fprintln(conn, "Login successful!")
		}
	}

}

// Keep track of the client's request time and only keep requests that are 10 seconds or less.
func updateRequestTracker(clientAddr string) {
	now := time.Now()
	timestamps := requestTracker[clientAddr]

	// Only keep timestamps within 10 seconds
	var validTimestamps []time.Time
	for _, t := range timestamps {
		if now.Sub(t) <= 10*time.Second {
			validTimestamps = append(validTimestamps, t)
		}
	}

	// Add the current time
	validTimestamps = append(validTimestamps, now)
	requestTracker[clientAddr] = validTimestamps
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
	logFile, err := os.OpenFile("yseo8834_jkim9115_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(logFile)
	} else {
		log.Fatal("Failed to log to file, using default stderr")
	}

	// Intended to set log output to logFile.
	log.SetOutput(logFile)

	// UDP > net.Listenpacket, TCP > net.listen
	listener, err := net.Listen("tcp", ":8080") // If an error occurs in this action, it is stored in the err variable.

	fmt.Println("Program Start!")

	// nil means 'nothing' or 'no error' in Go
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Info("Server Started. Listening on port 8080")

	clientID := 0 // Create for client-specific separation
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("Connection acceptance errors:", err)
			continue
		}

		clientID++
		go handleConnection(conn, clientID)
	}

}
