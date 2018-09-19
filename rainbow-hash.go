package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Entry struct {
	Preimage string `json:"preimage"`
	Md5      string `json:"md5"`
	Sha1     string `json:"sha1"`
	Sha256   string `json:"sha256"`
	Sha512   string `json:"sha512"`
}

var fileMutex sync.Mutex

func md5sum(preimage string) string {
	digest := md5.New()
	io.WriteString(digest, preimage)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func sha1sum(preimage string) string {
	digest := sha1.New()
	io.WriteString(digest, preimage)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func sha256sum(preimage string) string {
	digest := sha256.New()
	io.WriteString(digest, preimage)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func sha512sum(preimage string) string {
	digest := sha512.New()
	io.WriteString(digest, preimage)
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func computeWord(word string) ([]byte, error) {
	var entry Entry
	preimage := strings.TrimSpace(word)

	entry.Preimage = preimage
	entry.Md5 = md5sum(preimage)
	entry.Sha1 = sha1sum(preimage)
	entry.Sha256 = sha256sum(preimage)
	entry.Sha512 = sha512sum(preimage)
	return json.Marshal(entry)
}

func computeEntry(word string, output *os.File, wg *sync.WaitGroup) {
	data, err := computeWord(word)
	if err == nil {
		fileMutex.Lock()
		output.Write(data)
		output.Write([]byte("\n"))
		fileMutex.Unlock()
	}
	wg.Done()
}

func computeFile(input string, output string) {
	fInput, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer fInput.Close()

	fOutput, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer fOutput.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(fInput)
	counter := 0
	for scanner.Scan() {
		word := scanner.Text()
		wg.Add(1)
		go computeEntry(word, fOutput, &wg)
		counter++
		fmt.Printf("\rGo compute: %d", counter)
	}
	fmt.Printf("\rWaiting for go routines to finish ... ")
	wg.Wait()
	fmt.Printf("done")
}

func main() {
	fmt.Printf(" In: %s\n", os.Args[1])
	fmt.Printf("Out: %s\n", os.Args[2])
	computeFile(os.Args[1], os.Args[2])
}