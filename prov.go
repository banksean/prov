package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const BUFFSIZE = 1024

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <cmd> [args...] [> <output file>]\n", os.Args[0])
		os.Exit(1)
	}

	subCmd := os.Args[1]
	subArgs := os.Args[2:]

	cmd := exec.Command(subCmd, subArgs...)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %s\nError: %s\n", strings.Join(os.Args[2:], " "), err)
		os.Exit(1)
	}

	hash := sha1.New()
	hashWriter := bufio.NewWriter(hash)

	inbuf := make([]byte, BUFFSIZE)
	bread := 0

	// Read stdout from subprocess into buffer inbuf.
	// Then write inbuf to both the hasher and to the output file.

	reader := bufio.NewReader(&cmdOut)
	for ; err != io.EOF; bread, err = reader.Read(inbuf) {
		outbuf := bytes.NewReader(inbuf[:bread])
		if _, err = io.Copy(hashWriter, outbuf); err != nil {
			log.Fatal(err)
		}

		// Rewind so we can copy it back out to our stdout too.
		if _, err = outbuf.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		if _, err = io.Copy(os.Stdout, outbuf); err != nil {
			log.Fatal(err)
		}
	}

	hashWriter.Flush()

	// Write the provenance record for later query by whence.

	pFile, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".prov"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer pFile.Close()
	if err != nil {
		log.Fatalf("Couldn't open .prov: %v", err)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	line := []string{
		fmt.Sprintf("%x", hash.Sum(nil)),
		fmt.Sprintf("%v", time.Now()),
		os.Getenv("USER"),
		dir,
		fmt.Sprintf("%s", subCmd),
	}
	line = append(line, subArgs...)
	w := csv.NewWriter(pFile)

	err = w.Write(line)
	if err != nil {
		log.Fatalf("Couldn't write to .prov: %v", err)
	}
	w.Flush()
}
