package main

import (
	"fmt"
	
	"github.com/SatwikArnav/redis/store"

	"github.com/SatwikArnav/redis/RESP"
)

var Handler map[string]func([]RESP.Data) RESP.Data

func init() {
	Handler = make(map[string]func([]RESP.Data) RESP.Data)
	Handler["PING"] = PING
	Handler["COMMAND"] = COMMAND
	Handler["SET"] = SET
	Handler["GET"] = GET
	Handler["HSET"] = HSET
	Handler["HGET"] = HGET
}

func PING(D []RESP.Data) RESP.Data {
	if len(D) == 0 {
		return RESP.Data{
			CmdType: "+",
			Length:  4,
			Data:    "PONG"}
	}

	switch v := D[0].Data.(type) {
	case string:
		return RESP.Data{
			CmdType: "+",
			Length:  len(v) + 2,
			Data:    v,
		}
	case []byte:
		return RESP.Data{
			CmdType: "+",
			Length:  len(v) + 2,
			Data:    string(v),
		}
	case nil:
		return RESP.Data{
			CmdType: "+",
			Length:  5,
			Data:    "PONG",
		}

	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR wrong number of arguments for 'ping' command",
		}
	}
}

func COMMAND(D []RESP.Data) RESP.Data {
	// Return a simple OK response for now
	// In a full Redis implementation, this would return command documentation
	return RESP.Data{
		CmdType: "+",
		Length:  2,
		Data:    "OK",
	}
}

func SET(D []RESP.Data) RESP.Data {
	if len(D) != 2 {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR wrong number of arguments for 'set' command",
		}
	}
	
	// Convert key to string, handling both string and []byte types
	var key string
	switch v := D[0].Data.(type) {
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid key type",
		}
	}
	
	// Convert value to string, handling both string and []byte types
	var value string
	switch v := D[1].Data.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid value type",
		}
	}
	
	store.Set[key] = value
	fmt.Printf("Set key: %s, value: %s\n", key, value)
	fmt.Printf("Current store: %+v\n", store.Set)		
	return RESP.Data{
		CmdType: "+",
		Length:  2,
		Data:    "OK",
	}
}

func GET(D []RESP.Data) RESP.Data {
	if len(D) != 1 {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR wrong number of arguments for 'get' command",
		}
	}

	var key string
	switch v := D[0].Data.(type) {
	case string:
		key = v
	case []byte:
		key = string(v)
	default:
		return RESP.Data{	
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid key type",
		}
	}
	
	fmt.Printf("Getting key: %s\n", key)
	fmt.Printf("Current store: %+v\n", store.Set)
	value, exists := store.Set[key]
	if !exists {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR key not found",
		}
	}

	return RESP.Data{
		CmdType: "+",
		Length:  len(value) + 2,
		Data:    value,
	}
}

func HSET(D []RESP.Data) RESP.Data {
	if len(D) != 3 {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR wrong number of arguments for 'hset' command",
		}
	}
	var Pkey string
	switch v := D[0].Data.(type) {
	case string:
		Pkey = v	
	case []byte:
		Pkey = string(v)	
	default:
		return RESP.Data{	
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid key type",

		}
	}
	var field string
	switch v := D[1].Data.(type) {
	case string:
		field = v
	case []byte:
		field = string(v)
	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid field type",
		}
	}
	var value string
	switch v := D[2].Data.(type) {
	case string:
		value = v
	case []byte:
		value = string(v)
	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid value type",
		}
	}
	
	if _, exists := store.Hset[Pkey]; !exists {
		store.Hset[Pkey] = make(map[string]string)
	}
	store.Hset[Pkey][field] = value
	fmt.Printf("HSET key: %s, field: %s, value: %s\n", Pkey, field, value)
	fmt.Printf("Current HSet store: %+v\n", store.Hset)
	return RESP.Data{
		CmdType: "+",	
		Length:  2,
		Data:    "OK",
	}
}

func HGET(D []RESP.Data) RESP.Data {
	if len(D) != 2 {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR wrong number of arguments for 'hget' command",
		}
	}
	var Pkey string
	switch v := D[0].Data.(type) {
	case string:
		Pkey = v	
	case []byte:
		Pkey = string(v)	
	default:
		return RESP.Data{	
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid key type",
		}
	}
	var field string
	switch v := D[1].Data.(type) {
	case string:
		field = v
	case []byte:
		field = string(v)
	default:
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR invalid field type",
		}
	}

	fmt.Printf("HGET key: %s, field: %s\n", Pkey, field)
	fmt.Printf("Current HSet store: %+v\n", store.Hset)
	fields, exists := store.Hset[Pkey]
	if !exists {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR key not found",
		}
	}

	value, exists := fields[field]
	if !exists {
		return RESP.Data{
			CmdType: "-",
			Length:  0,
			Data:    "ERR field not found",
		}
	}

	return RESP.Data{
		CmdType: "+",
		Length:  len(value) + 2,
		Data:    value,
	}
}

