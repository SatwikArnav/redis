package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
	//defer client.Close()
	fmt.Println("Client connected:", client.RemoteAddr())

	reader := bufio.NewReader(client)
	
	// Keep connection alive and process multiple commands
	for {
		data, err := Read(reader)
		if err != nil {
			log.Println("Error reading data:", err)
			return
		}

		fmt.Printf("Received data: %+v\n", data)
		if data.cmdType != "*" {
			log.Println("NOT ARRAY")
			return
		}
		fmt.Println("data.data", data.data)
		fmt.Printf("data.data type: %T\n", data.data)
		var dataSlice []Data
		switch v := data.data.(type) {
		case []Data:
			dataSlice = v
		case []interface{}:
			for _, item := range v {
				d, ok := item.(Data)
				if !ok {
					log.Println("item is not of type Data")
					return
				}
				dataSlice = append(dataSlice, d)
			}
		default:
			log.Printf("data.data is not a []Data or []interface{}, but %T\n", data.data)
			return
		}
		fmt.Printf("dataSlice: %+v\n", dataSlice)
		fmt.Printf("dataSlice type: %T\n", dataSlice)
		if len(dataSlice) == 0 || data.length == 0 {
			log.Printf("dataSlice is empty or data.length is 0. dataSlice: %+v\n", dataSlice)
			return
		}
		command := dataSlice[0]
		args := dataSlice[1:]
		fmt.Printf("Command: %+v\n", command)
		fmt.Printf("Args: %+v\n", args)
		var cmdStr string
		switch v := command.data.(type) {
		case string:
			cmdStr = v
		case []byte:
			cmdStr = string(v)
		default:
			log.Printf("command.data is not a string or []byte, but %T\n", command.data)
			return
		}
		cmdStr = strings.ToUpper(cmdStr)
		fn, ok := Handler[cmdStr]
		if !ok {
			log.Println("Unknown command:", cmdStr)
			return
		}
		result := fn(args)

		fmt.Printf("Received data: %+v\n", data)

		response := marshall(result)
		fmt.Println("Response:", response)
		_, err = client.Write([]byte(response))
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
		fmt.Println("Response sent to client")
	}
}
