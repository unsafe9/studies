package main

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	googleAppsScriptWebAppURL = "<your-url>"
	zipName                   = "sheets.zip"
)

func main() {
	res, err := http.Get(googleAppsScriptWebAppURL)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	maxDecodedLen := base64.StdEncoding.DecodedLen(len(resBody))
	zip := make([]byte, maxDecodedLen)
	if n, err := base64.StdEncoding.Decode(zip, resBody); err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	} else {
		zip = zip[:n]
	}

	f, err := os.Create(zipName)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(zip)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	log.Printf("file saved as %s", zipName)
}
