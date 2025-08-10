package main

import (
	"fmt"
	//"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
)

func marshall(sata Data) string {
	//return fmt.Sprintf("hello%s%d\r\n\r\n", sata.cmdType, sata.length)
	DType := sata.cmdType
	//fmt.Println(DType)
	switch DType {
	case "+":
		return fmt.Sprintf("+%s\r\n", sata.data)
	case "-":
		return fmt.Sprintf("-%s\r\n", sata.data)

	case "$":
		return fmt.Sprintf("$%d\r\n%s\r\n", sata.length, sata.data)
	case "*":
		arr := sata.data.([]any)
		result := ""
		for i, item := range arr {
			if i == 0 {
				result += fmt.Sprintf("*%d\r\n%s", len(arr), marshall(item.(Data)))
			} else {
				result += fmt.Sprintf("%s%s", marshall(item.(Data)), "\r\n")
			}
		}
		return result
	case ":":
		return fmt.Sprintf(":%d\r\n", sata.data)

	default:
		return fmt.Sprintf("Unknown type: %s", sata.cmdType)

	}
}
