package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gqgs/flactomp3/pkg/notify"
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

	var notifier notify.Notifier
	notifier = notify.NewNotifier("flactomp3", filepath.Base(o.input))
	notifier.Start()
	if err := process(o); err != nil {
		notifier.Error(err)
		log.Fatal(err)
	}
	notifier.Success()
}
