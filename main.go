package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gen2brain/beeep"
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

	base := filepath.Base(o.input)

	beeep.Notify("flactomp3", fmt.Sprintf("Processing: %s", base), "/usr/share/icons/gnome/32x32/emblems/emblem-documents.png")
	if err := process(o); err != nil {
		beeep.Notify("flactomp3", fmt.Sprintf("Error: %s (%q)", base, err), "/usr/share/icons/gnome/32x32/emblems/emblem-important.png")
		log.Fatal(err)
	}

	beeep.Notify("flactomp3", fmt.Sprintf("Done (%s)", base), "/usr/share/icons/gnome/32x32/emblems/emblem-default.png")

}
