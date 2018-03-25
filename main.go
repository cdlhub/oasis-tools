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

	log.Printf("Write %q\n", rpFileName)
	log.Printf(" - cut index:         %d\n", j)
	log.Printf(" - cut return period: %d\n", m)

	fmt.Println("return_period")
	line := 1
	var rp int32
	for rp = int32(options.min); rp < m; rp += s {
		fmt.Println(rp)
		err := writeBin(f, rp)
		if err != nil {
			logCannotWrite(rpFileName, line, err)
		}
		line++
	}

	for i := j; rp < n; i-- {
		last := rp
		rp = n / i
		if rp < last+s/2 {
			continue
		}
		fmt.Println(rp)
		err := writeBin(f, rp)
		if err != nil {
			logCannotWrite(rpFileName, line, err)
		}
		line++
	}
}
