// This code is for the server side.

package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	fmt.Println("Hello, World!")

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
		fmt.Println("Check [4]")
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connected from client %v\n", conn.RemoteAddr())
	// 여기에 메시지 읽기 및 로깅 코드를 추가할 것!!!
}
