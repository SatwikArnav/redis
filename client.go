package main

import "net"
//tcp client

func SendMessage(message []byte){
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	println(string(buffer[:n]))

}