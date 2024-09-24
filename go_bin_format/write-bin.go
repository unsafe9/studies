package main

import (
	"encoding/binary"
	"encoding/json"
	"log"
	"os"
)

func main() {
	testData := map[string]any{
		"int":    1,
		"string": "hello",
		"bool":   true,
		"float":  1.1,
		"slice":  []int{1, 2, 3},
		"map":    map[string]int{"a": 1, "b": 2},
	}
	jsonData, _ := json.MarshalIndent(testData, "", "  ")
	log.Printf("json data: %s", jsonData)

	f, err := os.Create("test.bin")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(jsonData)))

	_, err = f.Write(header)
	if err != nil {
		log.Fatalf("failed to write header: %v", err)
	}

	_, err = f.Write(jsonData)
	if err != nil {
		log.Fatalf("failed to write data: %v", err)
	}

	log.Println("done")
}
