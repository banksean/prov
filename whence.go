package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const BUFFSIZE = 1024

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <filename>", os.Args[0])
		os.Exit(1)
	}

	pFile, err := os.Open(filepath.Join(os.Getenv("HOME"), ".prov"))
	if err != nil {
		log.Fatalf("Error opening .prov file: %v", err)
	}

	infile := os.Args[1]
	in, err := os.Open(infile)
	if err != nil {
		log.Fatalf("Error opening infile: %v", err)
	}
	inReader := bufio.NewReader(in)

	h := sha1.New()
	hWriter := bufio.NewWriter(h)
	buff := make([]byte, BUFFSIZE)
	bread := 0
	for ; err != io.EOF; bread, err = inReader.Read(buff) {
		_, err := hWriter.Write(buff[:bread])
		if err != nil {
			log.Printf("Error writing to sha1 hasher: %v", err)
		}
	}

	hWriter.Flush()
	cs := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("Hash: %s\n", cs)

	r := csv.NewReader(pFile)
	line := []string{}
	err = nil
	for ; err != io.EOF; line, err = r.Read() {
		if len(line) < 5 {
			continue
		}
		if line[0] == cs {
			fmt.Printf("Time: %v\n", line[1])
			fmt.Printf("User: %v\n", line[2])
			fmt.Printf("Directory: %v\n", line[3])
			fmt.Printf("Command: %v\n", strings.Join(line[4:], " "))
			os.Exit(0)
		}
	}

	log.Fatalf("Unrecognized content hash.")
}
