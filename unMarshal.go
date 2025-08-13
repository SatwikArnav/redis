package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Data struct {
	cmdType string
	length  int
	data    any
}

func readLength(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return 0, err
	}
	fmt.Println(string(line), "line")

	s := strings.TrimSpace(string(line))
	return strconv.Atoi(s)
}

func Read(reader *bufio.Reader) (Data, error) {
	typ, err := reader.ReadByte()
	if err != nil {
		return Data{}, err // Return error instead of calling log.Fatal
	}
	data := Data{
		cmdType: string(typ),
		length:  0,
		data:    nil,
	}
	switch typ {
	case '+':
		err = ReadString(reader, &data)
	case '-':
		err = ReadError(reader, &data)
	case '$':
		err = ReadBulk(reader, &data)
	case '*':
		err = ReadArray(reader, &data)
	case ':':
		err = ReadInteger(reader, &data)
	}
	if err != nil {
		return Data{}, err
	}
	return data, nil
}
func ReadString(reader *bufio.Reader, data *Data) error {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	data.length = len(line) - 2 // Exclude CRLF
	data.data = strings.TrimSpace(string(line))
	return nil
}

func ReadError(reader *bufio.Reader, data *Data) error {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	data.length = len(line) - 2 // Exclude CRLF
	data.data = strings.TrimSpace(string(line))
	return nil
}

func ReadBulk(reader *bufio.Reader, data *Data) error {
	length, err := readLength(reader)
	if err != nil {
		return err
	}
	data.length = length

	buf := make([]byte, length)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return err
	}
	data.data = buf

	// consume trailing CRLF
	if _, err = reader.ReadByte(); err != nil {
		return err
	}
	if _, err = reader.ReadByte(); err != nil {
		return err
	}
	return nil
}

func ReadInteger(reader *bufio.Reader, data *Data) error {
	value, err := readLength(reader)
	if err != nil {
		return err
	}
	data.data = value
	return nil
}

func ReadArray(reader *bufio.Reader, data *Data) error {
	length, err := readLength(reader)
	if err != nil {
		return err
	}
	data.length = length

	arr := make([]any, length)
	for i := range length {
		item, err := Read(reader)
		if err != nil {
			return err
		}
		arr[i] = item
	}
	data.data = arr
	return nil
}
