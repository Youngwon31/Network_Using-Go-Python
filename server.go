// This code is for the server side.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// Set up a log file and include log...
	// study:
	// The os.OpenFile function is the standard Go function for opening files.
	// os.O_CREATE: This flag creates the file if it doesn't exist.
	// os.O_WRONLY:  This flag is for opening files in write-only mode

	logFile, err := os.OpenFile("Yseo8834_jiukim9115_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err) // log.Fatal writes a log message and then calls os.Exit(1) to exit the program
	}
	defer logFile.Close() // The defer keyword schedules a specific function to run immediately before the current function (main or another function) returns.

	// Intended to set log output to logFile.
	log.SetOutput(logFile)

	fmt.Println("Hello, World!")

	// UDP > net.Listenpacket, TCP > net.listen
	listener, err := net.Listen("tcp", ":8080") // If an error occurs in this action, it is stored in the err variable.

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
	log.Printf("Connected from client %v\n", conn.RemoteAddr())
	// 여기에 메시지 읽기 및 로깅 코드를 추가할 것!!!
}
