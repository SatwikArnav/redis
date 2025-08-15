package RESP

import (
	"fmt"
	//"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
)

func Marshall(data Data) string {
	DType := data.CmdType
	switch DType {
	case "+":
		return fmt.Sprintf("+%s\r\n", data.Data)
	case "-":
		return fmt.Sprintf("-%s\r\n", data.Data)
	case "$":
		return fmt.Sprintf("$%d\r\n%s\r\n", data.Length, data.Data)
	case "*":
		switch arr := data.Data.(type) {
		case []Data:
			result := fmt.Sprintf("*%d\r\n", len(arr))
			for _, item := range arr {
				result += Marshall(item)
			}
			return result
		case []interface{}:
			result := fmt.Sprintf("*%d\r\n", len(arr))
			for _, item := range arr {
				if dataItem, ok := item.(Data); ok {
					result += Marshall(dataItem)
				}
			}
			return result
		default:
			return fmt.Sprintf("-ERR invalid array type: %T\r\n", data.Data)
		}
	case ":":
		return fmt.Sprintf(":%d\r\n", data.Data)
	default:
		return fmt.Sprintf("-ERR unknown type: %s\r\n", data.CmdType)
	}
}
