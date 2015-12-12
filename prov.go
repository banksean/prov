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
	"time"
)

const BUFFSIZE = 1024

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s filename <cmd> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	outfile := os.Args[1]
	out, err := os.Create(outfile)
	defer out.Close()
	if err != nil {
		log.Fatalf("Couldn't create %s: %v", outfile, err)
	}

	subCmd := os.Args[2]
	subArgs := os.Args[3:]

	cmd := exec.Command(subCmd, subArgs...)

	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut

	reader := bufio.NewReader(&cmdOut)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	hash := sha1.New()
	hashWriter := bufio.NewWriter(hash)

	inbuf := make([]byte, BUFFSIZE)
	bread := 0

	// Read stdout from subprocess into buffer inbuf.
	// Then write inbuf to both the hasher and to the output file.

	for ; err != io.EOF; bread, err = reader.Read(inbuf) {
		hashbuf := bytes.NewBuffer(inbuf[:bread])
		if _, err := io.Copy(hashWriter, hashbuf); err != nil {
			log.Fatalf(err.Error())
		}
		outbuf := bytes.NewBuffer(inbuf[:bread])
		if _, err := io.Copy(out, outbuf); err != nil {
			log.Fatalf(err.Error())
		}
	}

	hashWriter.Flush()

	pFile, err := os.OpenFile(".prov", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer pFile.Close()
	if err != nil {
		log.Fatalf("Couldn't open .prov: %v", err)
	}

	line := []string{
		fmt.Sprintf("%x", hash.Sum(nil)),
		fmt.Sprintf("%v", time.Now()),
		os.Getenv("USER"),
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
