package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gqgs/flactomp3/pkg/convert"
)

func updateFilename(path, format string) string {
	if !strings.Contains(strings.ToLower(path), "flac") {
		return fmt.Sprintf("%s [%s]", path, format)
	}

	var builder strings.Builder
	var i int
	builder.Grow(len(path))
	for i = 0; i <= len(path)-4; i++ {
		substr := path[i : i+4]
		if strings.EqualFold(substr, "flac") {
			n, _ := builder.WriteString(format)
			i += n + 1
			continue
		}
		builder.WriteByte(path[i])
	}

	if i <= len(path) {
		builder.WriteString(path[i:])
	}

	return builder.String()
}

func process(o options) error {
	if o.input == "" {
		return errors.New("input must be defined")
	}

	converter, err := convert.NewConverter(o.converter)
	if err != nil {
		return err
	}

	baseFolder := filepath.Base(o.input)

	outFolder := filepath.Join(o.output, updateFilename(baseFolder, "V0"))
	if err := os.Mkdir(outFolder, os.ModePerm); err != nil {
		return err
	}

	var hasConvertibleFiles bool

	err = filepath.WalkDir(o.input, func(path string, info fs.DirEntry, err error) error {
		relativePath := strings.TrimPrefix(path, o.input)
		relativePath = strings.TrimLeft(relativePath, "/")
		if relativePath == "" {
			return nil
		}

		if info.IsDir() {
			newFolder := filepath.Join(outFolder, relativePath)
			return os.Mkdir(newFolder, os.ModePerm)
		}

		hasConvertibleFiles = hasConvertibleFiles || convert.IsConvertible(info.Name())

		return nil
	})

	if err != nil {
		return err
	}

	if !hasConvertibleFiles {
		originalName := filepath.Join(o.output, baseFolder)
		if err := os.Rename(outFolder, originalName); err != nil {
			return err
		}
		outFolder = originalName
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

	err = filepath.WalkDir(o.input, func(path string, info fs.DirEntry, err error) error {
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
