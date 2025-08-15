package RESP

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Data struct {
	CmdType string
	Length  int
	Data    any
}

func ReadLength(reader *bufio.Reader) (int, error) {
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
		CmdType: string(typ),
		Length:  0,
		Data:    nil,
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
	data.Length = len(line) - 2 // Exclude CRLF
	data.Data = strings.TrimSpace(string(line))
	return nil
}

func ReadError(reader *bufio.Reader, data *Data) error {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	data.Length = len(line) - 2 // Exclude CRLF
	data.Data = strings.TrimSpace(string(line))
	return nil
}

func ReadBulk(reader *bufio.Reader, data *Data) error {
	length, err := ReadLength(reader)
	if err != nil {
		return err
	}
	data.Length = length

	buf := make([]byte, length)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return err
	}
	data.Data = buf

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
	value, err := ReadLength(reader)
	if err != nil {
		return err
	}
	data.Data = value
	return nil
}

func ReadArray(reader *bufio.Reader, data *Data) error {
	length, err := ReadLength(reader)
	if err != nil {
		return err
	}
	data.Length = length

	arr := make([]any, length)
	for i := range length {
		item, err := Read(reader)
		if err != nil {
			return err
		}
		arr[i] = item
	}
	data.Data = arr
	return nil
}
