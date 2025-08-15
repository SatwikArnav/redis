package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"github.com/SatwikArnav/redis/RESP"
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
	defer client.Close() // Uncomment this to properly close connections
	fmt.Println("Client connected:", client.RemoteAddr())

	reader := bufio.NewReader(client)

	// Keep connection alive and process multiple commands
	for {
		data, err := RESP.Read(reader)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client disconnected:", client.RemoteAddr())
			} else {
				log.Println("Error reading data:", err)
			}
			return
		}

		fmt.Printf("Received data: %+v\n", data)
		if data.CmdType != "*" {
			log.Println("NOT ARRAY")
			client.Write([]byte("-ERR invalid command format\r\n"))
			continue
		}
		fmt.Println("data.data", data.Data)
		fmt.Printf("data.data type: %T\n", data.Data)
		var dataSlice []RESP.Data
		switch v := data.Data.(type) {
		case []RESP.Data:
			dataSlice = v
		case []interface{}:
			for _, item := range v {
				d, ok := item.(RESP.Data)
				if !ok {
					log.Println("item is not of type Data")
					client.Write([]byte("-ERR invalid data format\r\n"))
					continue
				}
				dataSlice = append(dataSlice, d)
			}
		default:
			log.Printf("data.data is not a []Data or []interface{}, but %T\n", data.Data)
			client.Write([]byte("-ERR invalid data format\r\n"))
			continue
		}
		fmt.Printf("dataSlice: %+v\n", dataSlice)
		fmt.Printf("dataSlice type: %T\n", dataSlice)
		if len(dataSlice) == 0 || data.Length == 0 {
			log.Printf("dataSlice is empty or data.length is 0. dataSlice: %+v\n", dataSlice)
			client.Write([]byte("-ERR empty command\r\n"))
			continue
		}
		command := dataSlice[0]
		args := dataSlice[1:]
		fmt.Printf("Command: %+v\n", command)
		fmt.Printf("Args: %+v\n", args)
		var cmdStr string
		switch v := command.Data.(type) {
		case string:
			cmdStr = v
		case []byte:
			cmdStr = string(v)
		default:
			log.Printf("command.data is not a string or []byte, but %T\n", command.Data)
			client.Write([]byte("-ERR invalid command format\r\n"))
			continue
		}
		cmdStr = strings.ToUpper(cmdStr)
		fn, ok := Handler[cmdStr]
		if !ok {
			log.Println("Unknown command:", cmdStr)
			client.Write([]byte("-ERR unknown command '" + cmdStr + "'\r\n"))
			continue
		}
		result := fn(args)

		fmt.Printf("Received data: %+v\n", data)

		response := RESP.Marshall(result)
		fmt.Println("Response:", response)
		_, err = client.Write([]byte(response))
		if err != nil {
			log.Println("Error writing response:", err)
			return
		}
		fmt.Println("Response sent to client")
	}
}
