package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var num int

func init() {
	flag.IntVar(&num, "n", 1, "number of lines to sample")
}

func main() {
	flag.Parse()

	var lines []string

	scan := bufio.NewScanner(os.Stdin)
	for i := 0; i < num; i++ {
		if !scan.Scan() {
			log.Fatal("Too few lines in input")
		}

		lines = append(lines, scan.Text())
	}

	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	seen := num
	for scan.Scan() {
		seen++
		index := r.Intn(seen)
		if index < num {
			lines[index] = scan.Text()
		}
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}
