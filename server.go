package main

import (
	"bufio"
	"net"

	"fmt"
	"log"
)

func TcpListener() {
	conn, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Server is listening on port 6379...")
	for {
		client, err := conn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleClient(client)
	}
}

func HandleClient(client net.Conn) {
	defer client.Close()
	fmt.Println("Client connected:", client.RemoteAddr())

	reader := bufio.NewReader(client)
	data, err := Read(reader)
	if err != nil {
		log.Println("Error reading data:", err)
		return
	}
	fmt.Printf("Received data: %+v\n", data)

	response := marshall(data)
	_, err = client.Write([]byte(response))
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}
	fmt.Println("Response sent to client")
}
