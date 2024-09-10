package main

import (
	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

//go:generate go run $GOFILE

func main() {
	rangeSchemaFiles(func(file string) {
		replacePSQLCommands(file, true)
	})

	exitCode := sqlc.Run([]string{"generate"})

	rangeSchemaFiles(func(file string) {
		replacePSQLCommands(file, false)
	})

	os.Exit(exitCode)
}

func rangeSchemaFiles(callback func(file string)) {
	sqlFiles, err := filepath.Glob("../schemas/*.sql")
	if err != nil {
		log.Fatalf("failed to list schema file: %v", err)
	}

	wg := sync.WaitGroup{}
	for _, file := range sqlFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			callback(file)
		}(file)
	}
	wg.Wait()
}

func replacePSQLCommands(file string, commentOut bool) {
	var (
		re   *regexp.Regexp
		into []byte
	)
	if commentOut {
		re = regexp.MustCompile(`(?m)^\\`)
		into = []byte("--<sqlc.go> \\")
	} else {
		re = regexp.MustCompile(`(?m)^--<sqlc.go> \\`)
		into = []byte("\\")
	}

	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("failed to read %s: %v", file, err)
	}

	modified := re.ReplaceAll(content, into)

	err = os.WriteFile(file, modified, 0644)
	if err != nil {
		log.Fatalf("failed to write %s: %v", file, err)
	}
}
