package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type options struct {
	output   string
	input    string
	parallel int
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var o options
	flag.StringVar(&o.input, "input", "", "input folder")
	flag.StringVar(&o.output, "output", wd, "output folder")
	flag.IntVar(&o.parallel, "parallel", runtime.NumCPU(), "number of parallel processes")
	flag.Parse()

	if err := process(o); err != nil {
		log.Fatal(err)
	}

}

func process(o options) error {
	baseFolder := filepath.Base(o.input)
	log.Printf("Processing: %q", baseFolder)

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

					if isConvertible(relativePath) {
						return convert(relativePath, o.input, outFolder)
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

func isConvertible(path string) bool {
	return filepath.Ext(path) == ".flac"
}

func convert(relativePath, baseFolder, outFolder string) error {
	mp3FileName := strings.TrimSuffix(relativePath, ".flac") + ".mp3"
	newPath := filepath.Join(outFolder, mp3FileName)
	originalPath := filepath.Join(baseFolder, relativePath)
	log.Printf("Creating: %q", mp3FileName)
	cmd := exec.Command("ffmpeg", "-i", originalPath, "-q:a", "0", "-threads", "0", "-loglevel", "error", newPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
