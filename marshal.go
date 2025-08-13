package main

import (
	"fmt"
	//"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
)

func marshall(data Data) string {
	DType := data.cmdType
	switch DType {
	case "+":
		return fmt.Sprintf("+%s\r\n", data.data)
	case "-":
		return fmt.Sprintf("-%s\r\n", data.data)
	case "$":
		return fmt.Sprintf("$%d\r\n%s\r\n", data.length, data.data)
	case "*":
		switch arr := data.data.(type) {
		case []Data:
			result := fmt.Sprintf("*%d\r\n", len(arr))
			for _, item := range arr {
				result += marshall(item)
			}
			return result
		case []interface{}:
			result := fmt.Sprintf("*%d\r\n", len(arr))
			for _, item := range arr {
				if dataItem, ok := item.(Data); ok {
					result += marshall(dataItem)
				}
			}
			return result
		default:
			return fmt.Sprintf("-ERR invalid array type: %T\r\n", data.data)
		}
	case ":":
		return fmt.Sprintf(":%d\r\n", data.data)
	default:
		return fmt.Sprintf("-ERR unknown type: %s\r\n", data.cmdType)
	}
}
