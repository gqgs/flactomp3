package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gqgs/flactomp3/pkg/notify"
)

//go:generate go tool argsgen

type options struct {
	input     string `arg:"input folder,positional"`
	output    string `arg:"output folder,positional"`
	converter string `arg:"output encoding format (opus | lame)"`
	parallel  int    `arg:"number of parallel processes"`
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	o := options{
		output:    wd,
		parallel:  runtime.NumCPU(),
		converter: "lame",
	}
	o.MustParse()

	notifier := notify.NewNotifier("flactomp3", filepath.Base(o.input))
	notifier.Start()
	if err := process(o); err != nil {
		notifier.Error(err)
		log.Fatal(err)
	}
	notifier.Success()
}
