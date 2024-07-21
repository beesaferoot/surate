package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

// loadFromDisk loads the rate data and checks the timestamp
func loadFromDisk(filename string) (RateResponse, error) {
	var data RateResponse
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	// Read the timestamp
	var timestamp int64
	n, err := fmt.Sscanf(string(fileData), "%d\n", &timestamp)
	if n != 1 || err != nil {
		return data, fmt.Errorf("failed to read timestamp")
	}

	// Check if the file is older than an hour
	if time.Now().Unix()-timestamp > 3600 {
		return data, fmt.Errorf("file is older than an hour")
	}

	// Decode the data
	buf := bytes.NewBuffer(fileData[len(fmt.Sprintf("%d\n", timestamp)):])
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&data)
	return data, err
}

func saveToDisk(data RateResponse, filename string) error {
	var buf bytes.Buffer

	// Add a timestamp to the beginning of the file
	timestamp := time.Now().Unix()
	buf.Write([]byte(fmt.Sprintf("%d\n", timestamp)))

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}
