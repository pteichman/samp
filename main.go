package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// samp performs stream sampling on the command line. It's the lossy cat(1)
// you never knew you needed.

func main() {
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	flag0 := fs.Bool("0", false, "use NUL characters instead of line delimiters")
	flagK := fs.Int("k", 1, "maximum number of items to pass")
	flagS := fs.Int64("s", -1, "random seed: -1 for current time")

	fs.Parse(os.Args[1:])

	scan := bufio.NewScanner(os.Stdin)
	if *flag0 {
		scan.Split(splitNulls)
	}

	seed := *flagS
	if seed < 0 {
		seed = time.Now().UnixNano()
	}
	r := rand.New(rand.NewSource(seed))

	// Fill the reservoir of lines with an initial K items.
	lines := make([]line, 0, *flagK)
	for i := 0; i < *flagK; i++ {
		if scan.Scan() {
			lines = append(lines, line{i, dup(scan.Bytes())})
		}
	}

	// Randomly replace items in the reservoir with the remaining lines.
	for count := len(lines); scan.Scan(); count++ {
		index := r.Intn(count)
		if index < len(lines) {
			lines[index] = line{count, dup(scan.Bytes())}
		}
	}

	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %s", err)
		os.Exit(1)
	}

	// Sort the reservoir to ensure output is ordered like the input.
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].num < lines[j].num
	})

	delim := '\n'
	if *flag0 {
		delim = 0
	}

	// Write the reservoir to stdout.
	w := bufio.NewWriter(os.Stdout)
	for _, line := range lines {
		_, err := w.Write(line.data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %s", err)
			os.Exit(1)
		}
		_, err = w.WriteRune(delim)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Write error: %s", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %s", err)
		os.Exit(1)
	}
}

type line struct {
	num  int
	data []byte
}

func splitNulls(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func dup(b []byte) []byte {
	ret := make([]byte, len(b))
	copy(ret, b)
	return ret
}
