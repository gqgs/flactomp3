package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gqgs/flactomp3/pkg/convert"
)

func process(o options) error {
	if o.input == "" {
		return errors.New("input must be defined")
	}

	var converter convert.Converter
	converter = convert.NewFFmpegConverter()

	baseFolder := filepath.Base(o.input)

	outFolder := filepath.Join(o.output, baseFolder+" [V0]")
	if err := os.Mkdir(outFolder, os.ModePerm); err != nil {
		return err
	}

	err := filepath.Walk(o.input, func(path string, info os.FileInfo, err error) error {
		relativePath := strings.TrimPrefix(path, o.input)
		relativePath = strings.TrimLeft(relativePath, "/")
		if relativePath == "" {
			return nil
		}

		if info.IsDir() {
			newFolder := filepath.Join(outFolder, relativePath)
			return os.Mkdir(newFolder, os.ModePerm)
		}

		return nil
	})

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	pathCh := make(chan string)

	go func() {
		semaphore := make(chan struct{}, o.parallel)
		for path := range pathCh {
			path := path
			semaphore <- struct{}{}
			go func() {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				if err := func() error {
					relativePath := strings.TrimPrefix(path, o.input)
					relativePath = strings.TrimLeft(relativePath, "/")

					if convert.IsConvertible(relativePath) {
						return converter.Convert(relativePath, o.input, outFolder)
					}

					inFile, err := os.Open(path)
					if err != nil {
						return err
					}
					defer inFile.Close()

					outFile, err := os.Create(filepath.Join(outFolder, relativePath))
					if err != nil {
						return err
					}
					defer outFile.Close()

					_, err = io.Copy(outFile, inFile)
					return err
				}(); err != nil {
					log.Printf("%s: %s", err, path)
				}

			}()
		}
	}()

	err = filepath.Walk(o.input, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		pathCh <- path

		return nil
	})
	wg.Wait()

	return err
}