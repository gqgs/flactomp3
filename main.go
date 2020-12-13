package main

import (
	"flag"
	"log"
	"os"
	"runtime"
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
