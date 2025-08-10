package main

import (
	"fmt"

	"github.com/SatwikArnav/redis/store"
)

var Handler map[string]func([]Data) Data

func init() {
	Handler = make(map[string]func([]Data) Data)
	Handler["PING"] = PING
	Handler["COMMAND"] = COMMAND
	Handler["SET"] = SET
	Handler["GET"] = GET
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

func SET(D []Data) Data {
	if len(D) != 2 {
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR wrong number of arguments for 'set' command",
		}
	}
	
	// Convert key to string, handling both string and []byte types
	var key string
	switch v := D[0].data.(type) {
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR invalid key type",
		}
	}
	
	// Convert value to string, handling both string and []byte types
	var value string
	switch v := D[1].data.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR invalid value type",
		}
	}
	
	store.Set[key] = value
	fmt.Printf("Set key: %s, value: %s\n", key, value)
	fmt.Printf("Current store: %+v\n", store.Set)		
	return Data{
		cmdType: "+",
		length:  2,
		data:    "OK",
	}
}

func GET(D []Data) Data {
	if len(D) != 1 {
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR wrong number of arguments for 'get' command",
		}
	}

	var key string
	switch v := D[0].data.(type) {
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		return Data{	
			cmdType: "-",
			length:  0,
			data:    "ERR invalid key type",
		}
	}
	
	fmt.Printf("Getting key: %s\n", key)
	fmt.Printf("Current store: %+v\n", store.Set)
	value, exists := store.Set[key]
	if !exists {
		return Data{
			cmdType: "-",
			length:  0,
			data:    "ERR key not found",
		}
	}

	return Data{
		cmdType: "+",
		length:  len(value) + 2,
		data:    value,
	}
}