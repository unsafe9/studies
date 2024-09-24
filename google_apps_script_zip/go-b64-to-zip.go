package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	googleAppsScriptWebAppURL = "<your-url>"
	zipName                   = "sheets.zip"
	dirName                   = "sheets"
)

func main() {
	zipData := downloadZip()
	writeZip(zipData)
	unzipToDir(zipData)
}

func downloadZip() []byte {
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
	zipData := make([]byte, maxDecodedLen)
	n, err := base64.StdEncoding.Decode(zipData, resBody)
	if err != nil {
		log.Fatalf("failed to decode base64: %v", err)
	}
	return zipData[:n]
}

func writeZip(zipData []byte) {
	f, err := os.Create(zipName)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(zipData)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	log.Printf("file saved as %s", zipName)
}

func unzipToDir(zipData []byte) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		log.Fatalf("failed to read zip: %v", err)
	}

	os.MkdirAll(dirName, os.ModePerm)
	wg := sync.WaitGroup{}
	wg.Add(len(zipReader.File))
	for _, f := range zipReader.File {
		go func() {
			defer wg.Done()

			rc, err := f.Open()
			if err != nil {
				log.Fatalf("failed to open file: %v", err)
			}
			defer rc.Close()

			outFile, err := os.Create(filepath.Join(dirName, f.Name))
			if err != nil {
				log.Fatalf("failed to create file: %v", err)
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				log.Fatalf("failed to copy file: %v", err)
			}
		}()
	}
	wg.Wait()
}
