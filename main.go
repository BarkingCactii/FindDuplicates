package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var start = time.Now()
	argsWithoutProg := os.Args[1]
	fmt.Println("Processing -> ", argsWithoutProg)
	walkdir(argsWithoutProg)
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %f seconds", elapsed.Seconds())
}

func walkdir(path string) error {
	var count int64
	var duplicates int64

	m := map[string][]string{}
	count = 0
	duplicates = 0

	var scan = func(path string, info os.FileInfo, err error) error {

		if info.IsDir() == true {
			return nil
		}

		count++
		if count%31 == 0 {
			fmt.Printf("Processed %d files\r", count)
		}

		hash := readfile(path)
		shaStr := hex.EncodeToString(hash)
		m[shaStr] = append(m[shaStr], path)
		return nil
	}

	err := filepath.Walk(path, scan)

	for k, v := range m {
		var numFiles = len(v)
		if numFiles > 1 {
			fmt.Println(k)
			for _, v2 := range v {
				fmt.Println("\t" + v2)
			}
			duplicates += int64(numFiles - 1)
		}
	}

	fmt.Println("Processed", count, "file(s) with", duplicates, "duplicates")
	return err
}

const filechunk = 65535 // we settle for 8KB

func readfile(filename string) []byte {
	file, err := os.Open(filename)

	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()

	filesize := info.Size()

	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return hash.Sum(nil)
}
