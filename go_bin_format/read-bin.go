package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("test.bin")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	jsonBytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	if len(jsonBytes) < 4 {
		log.Fatalf("file is too short: %d", len(jsonBytes))
	}

	jsonBytes = jsonBytes[4:]

	data := map[string]any{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %v", err)
	}

	log.Printf("bytes: %d", len(jsonBytes))
	log.Printf("json: %s", string(jsonBytes))
	log.Printf("data: %+v", data)
}
