package persistence

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/SatwikArnav/redis/RESP"
)

type AOF struct {
	File *os.File
	Rd   *bufio.Reader
	Mu   sync.Mutex
}

func NewAOF(filePath string) (*AOF, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening AOF file: %w", err)
	}

	resp := bufio.NewReader(file)
	aof := &AOF{
		File: file,
		Rd:   resp,
	}

	go func() {
		aof.Mu.Lock()
		defer aof.Mu.Unlock()
		aof.File.Sync()
		time.Sleep(time.Second * 5) // Sync every 5 seconds
	}()

	return aof, nil
}

func CloseAOF(aof *AOF) {
	if aof != nil && aof.File != nil {
		aof.Mu.Lock()
		defer aof.Mu.Unlock()
		if err := aof.File.Close(); err != nil {
			fmt.Printf("Error closing AOF file: %v\n", err)
		}
	}
}

func (aof *AOF) Write(data RESP.Data) error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	// Marshal the data to RESP format
	marshaledData := RESP.Marshall(data)

	// Write the marshaled data to the AOF file
	_, err := aof.File.WriteString(marshaledData)
	if err != nil {
		return fmt.Errorf("error writing to AOF file: %w", err)
	}

	return nil
}

func (aof *AOF) Read(processData func(data RESP.Data) error) error {
	aof.Mu.Lock()
	defer aof.Mu.Unlock()

	for {
		data, err := RESP.Read(aof.Rd)
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			return fmt.Errorf("error reading from AOF file: %w", err)
		}

		if err := processData(data); err != nil {
			return fmt.Errorf("error processing data: %w", err)
		}
	}

	return nil
}