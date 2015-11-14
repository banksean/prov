package main

import (
	"bufio"
	//"bytes"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	//"os/exec"
	//"time"
)

func main() {
	pFile, err := os.Open(".prov")
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
	p := make([]byte, 1024)
	n := 0
	for ; err != io.EOF; n, err = inReader.Read(p) {
		_, err := hWriter.Write(p[:n])
		if err != nil {
			log.Printf("Error writing to sha1 hasher: %v", err)
		}
	}

	hWriter.Flush()
	cs := fmt.Sprintf("%x", h.Sum(nil))

	fmt.Printf("hash: %s\n", cs)

	r := csv.NewReader(pFile)
	line := []string{}
	err = nil
	for ; err != io.EOF; line, err = r.Read() {
		if len(line) < 4 {
			continue
		}
		if line[0] == cs {
			fmt.Printf("%v", line)
			os.Exit(0)
		}
	}

	log.Fatalf("Unrecognized content hash.")
}
