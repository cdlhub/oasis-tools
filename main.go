package main

import (
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

const rpFileName = "returnperiods.bin"

// Options is for command line options
type Options struct {
	min     int
	step    int
	max     int
	verbose bool
}

var options Options

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "\t%s -version\n", NAME)
		fmt.Fprintf(os.Stderr, "\t%s -help\n", NAME)
		fmt.Fprintf(os.Stderr, "\t%s [ -min M ] [ -step S ] [ -manx N ] [ -verbose ]\n", NAME)
		fmt.Fprintln(os.Stderr)

		flag.PrintDefaults()
	}

	flag.IntVar(&options.min, "min", 5, "The minimum return period")
	flag.IntVar(&options.step, "step", 5, "The minimum difference between two return periods")
	flag.IntVar(&options.max, "max", 10000, "The maximum return period")
	flag.BoolVar(&options.verbose, "verbose", false, "Vebose mode to output csv period file on standard output")
	flag.Parse()
}

func writeBin(f *os.File, i int32) error {
	return binary.Write(f, binary.LittleEndian, i)
}

func logCannotWrite(fileName string, line int, err error) {
	log.Fatalf("ERROR: cannot write line %d to file %q: %v", line, fileName, err)
}

func logStart() {
	log.Println("Options:")
	log.Printf(" - min return period: %d\n", options.min)
	log.Printf(" - max return period: %d\n", options.max)
	log.Printf(" - minumum step:      %d\n", options.step)
}

// Sample:
// - min return period: 5
// - max return period: 100
// - minumum step:      5
//
// return_period
// 5
// 10
// 15
// 20
// 33
// 50
// 100
func main() {
	log.SetFlags(0) // No time
	logStart()

	f, err := os.Create(rpFileName)
	if err != nil {
		log.Fatalf("ERROR: cannot create %q: %v", rpFileName, err)
	}
	defer f.Close()

	// m: min return period before step by step
	// s: step
	// n: max return period
	//     1   sqrt(4n + (s+1))
	// m = - * ---------------- - 1
	//     2      sqrt(s + 1)
	n := int32(options.max)
	s := int32(options.step)
	r1 := float64(4*n + s)
	r2 := float64(s)
	j := int32(.5 * (math.Sqrt(r1)/math.Sqrt(r2) - 1)) // + .5 for int coinversion
	m := n / j

	if options.verbose {
		fmt.Println("return_period")
	}

	line := 1

	var rp int32
	for i := int32(1); i <= j; i++ {
		rp = n / i
		if options.verbose {
			fmt.Println(rp)
		}
		err := writeBin(f, rp)
		if err != nil {
			logCannotWrite(rpFileName, line, err)
		}
		line++
	}

	for rp = rp - s - ((rp - s) % s); rp >= int32(options.min); rp -= s {
		if options.verbose {
			fmt.Println(rp)
		}
		err := writeBin(f, rp)
		if err != nil {
			logCannotWrite(rpFileName, line, err)
		}
		line++
	}

	log.Printf("Write %q\n", rpFileName)
	log.Println(" - linear part:")
	log.Printf("    start:  %12d\n", options.min)
	log.Printf("    end:    %12d\n", m)
	log.Printf("    points: %12d\n", j)
	log.Printf(" - total number of points: %d", line-1)
}
