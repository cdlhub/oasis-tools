package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

const (
	// VERSION sementic number of the command.
	// It is the only place to change version string.
	VERSION string = "0.1.0"
	// NAME is the name of the command.
	NAME string = "returnperiod"
)

// Options is for command line options
type Options struct {
	min  int
	step int
	max  int
}

var options Options

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "\t%s -version\n", NAME)
		fmt.Fprintf(os.Stderr, "\t%s -help\n", NAME)
		fmt.Fprintf(os.Stderr, "\t%s -rp <N>\n", NAME)
		fmt.Fprintln(os.Stderr)

		flag.PrintDefaults()
	}

	flag.IntVar(&options.min, "min", 5, "The minimum return period")
	flag.IntVar(&options.step, "step", 5, "The minimum difference between two return periods")
	flag.IntVar(&options.max, "max", 10000, "The maximum return period")
	flag.Parse()
}

func writeBin(f *os.File, i int) error {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, i)
	_, err := f.Write(buf.Bytes())
	return err
}

func main() {
	// R: min return period before step by step
	// S: step
	// N: max return period
	// R = 1/2 * ((sqrt(4*N+S) / sqrt(S)) - 1)
	n := options.max
	s := options.step + 1
	r1 := float64(4*n + s)
	r2 := float64(s)
	j := int(.5*(math.Sqrt(r1)/math.Sqrt(r2)-1) + .5)
	middleRP := options.max / j

	log.SetFlags(0) // No time
	log.Println("Options:")
	log.Printf(" - min return period: %d\n", options.min)
	log.Printf(" - max return period: %d\n", options.max)
	log.Printf(" - minumum step:      %d\n", options.step)

	f, err := os.Create("returnperiods.bin")
	if err != nil {
		log.Fatalf("ERROR: cannot write to 'returnperiods.bin': %v", err)
	}
	defer f.Close()

	var rp int
	for rp = options.min; rp < middleRP; rp += options.step {
		err := writeBin(f, rp)
		if err != nil {
			log.Fatalf("ERROR: cannot write to file 'returnperiods.bin': %v", err)
		}
	}

	for i := j; rp < options.max; i-- {
		rp = options.max / i
		err := writeBin(f, rp)
		if err != nil {
			log.Fatalf("ERROR: cannot write to file 'returnperiods.bin': %v", err)
		}
	}
}
