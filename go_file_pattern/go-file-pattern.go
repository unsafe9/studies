package main

import (
	"github.com/bmatcuk/doublestar/v4"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"path/filepath"
	"time"
)

const pattern = "**/*.md"

func main() {
	if !doublestar.ValidatePattern(pattern) {
		log.Fatalf("invalid pattern: %s", pattern)
	}

	ch := make(chan string, 1000)
	done := make(chan struct{})
	go func() {
		defer close(done)
		var wg errgroup.Group
		for f := range ch {
			wg.Go(func() error {
				log.Printf("processing %s", f)
				time.Sleep(time.Millisecond * 100)
				return nil
			})
		}
		if err := wg.Wait(); err != nil {
			log.Fatalf("failed at least one: %v", err)
		}
	}()

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if match, _ := doublestar.Match(pattern, path); match {
			ch <- path
		}
		return nil
	})
	if err != nil {
		log.Fatalf("failed to walk: %v", err)
	}

	close(ch)
	<-done
}
