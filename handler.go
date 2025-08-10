package main

//import "bufio"
//import "fmt"

var Handler map[string]func([]Data) Data

func init() {
	Handler = make(map[string]func([]Data) Data)
	Handler["PING"] = PING
	Handler["COMMAND"] = COMMAND
}

func PING(D []Data) Data {
	if len(D) == 0 {
		return Data{
			cmdType: "+",
			length:  4,
			data:    "PONG"}
	}

	switch v := D[0].data.(type) {
	case string:
		return Data{
			cmdType: "+",
			length:  len(v) + 2,
			data:    v,
		}
	case []byte:
		return Data{
			cmdType: "+",
			length:  len(v) + 2,
			data:    string(v),
		}
	case nil:
		return Data{
			cmdType: "+",
			length:  5,
			data:    "PONG",
		}

	default:
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR wrong number of arguments for 'ping' command",
		}
	}
}

func COMMAND(D []Data) Data {
	return Data{
		cmdType: "+",
		length:  2,
		data:    "OK",
	}
}
