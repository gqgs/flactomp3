// Code generated by argsgen.
// DO NOT EDIT!
package main

import (
    "flag"
    "fmt"
    "os"
)

func (o *options) flagSet() *flag.FlagSet {
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&o.input, "input", o.input, "input folder")
    flagSet.StringVar(&o.output, "output", o.output, "output folder")
    flagSet.StringVar(&o.converter, "converter", o.converter, "output encoding format (opus | lame)")
    flagSet.IntVar(&o.parallel, "parallel", o.parallel, "number of parallel processes")
    return flagSet
}

// Parse parses the arguments in os.Args
func (o *options) Parse() error {
    flagSet := o.flagSet()
    
    var positional []string
    args := os.Args[1:]
    for len(args) > 0 {
        if err := flagSet.Parse(args); err != nil {
            return err
        }

        if remaining := flagSet.NArg(); remaining > 0 {
            posIndex := len(args) - remaining
            
            positional = append(positional, args[posIndex])
            args = args[posIndex+1:]
            continue
        }
        break
    }

    
    if len(positional) == 0 {
        return nil
    }
    if len(positional) > 0 {
        o.input = positional[0]
    }
    if len(positional) > 1 {
        o.output = positional[1]
    }
    return nil
}

// MustParse parses the arguments in os.Args or exists on error
func (o *options) MustParse() {
    if err := o.Parse(); err != nil {
        o.flagSet().PrintDefaults()
        fmt.Fprintln(os.Stderr)
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
