package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
)

const MAX_SIZE = 470000
const DESIRED_ERROR_RATE = 0.001
const IDENTIFIER = "CCBF"
const VERSION_NUMBER uint16 = 1

func main() {
	var buildFlag string
	var outputFlag string
	var readFlag string

	flag.StringVar(&buildFlag, "build", "", "specify a file that will be used to create the bloom filter")
	flag.StringVar(&outputFlag, "output", "./tmp/words.bf", "specify an output file for the bloom filter")
	flag.StringVar(&readFlag, "read", "./tmp/words.bf", "specify a bloom filter file")
	flag.Parse()

	if buildFlag != "" {
		bloomFilter := buildBloomFilter(buildFlag, MAX_SIZE, DESIRED_ERROR_RATE)
		storeBloomFilter(bloomFilter, outputFlag)
	} else {
		startIndex := 1
		if readFlag != "./tmp/words.bf" {
			startIndex = 3
		}
		checkWordsSpelling(os.Args[startIndex:], readFlag)
	}
}

func buildBloomFilter(filePath string, maxSize float64, desiredErrorRate float64) *BloomFilter {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	bloomFilter := createBloomFilter(maxSize, desiredErrorRate)

	for {
		buffer, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		word := string(buffer)
		bloomFilter.insert(word)
	}

	return bloomFilter
}

func storeBloomFilter(bloomFilter *BloomFilter, outputFilePath string) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	identifier := []byte(IDENTIFIER)
	var versionNumber [2]byte
	var numHashApplied [2]byte
	var totalBits [8]byte

	binary.BigEndian.PutUint16(versionNumber[:], VERSION_NUMBER)
	binary.BigEndian.PutUint16(numHashApplied[:], uint16(bloomFilter.numHashApplied))
	binary.BigEndian.PutUint64(totalBits[:], bloomFilter.totalBits)

	writer.Write(identifier)
	writer.Write(versionNumber[:])
	writer.Write(numHashApplied[:])
	writer.Write(totalBits[:])
	writer.Write(bloomFilter.bloomFilter)

	writer.Flush()
}

func checkWordsSpelling(words []string, filePath string) {
	if len(words) < 1 {
		return
	}

	bloomFilter := readBloomFilter(filePath)
	wrongWords := make([]string, 0)

	for _, word := range words {
		if !bloomFilter.query(word) {
			wrongWords = append(wrongWords, word)
		}
	}

	if len(wrongWords) < 1 {
		fmt.Println("No words were spelt wrong!")
	} else {
		fmt.Println("These words were spelt wrong:")

		for _, word := range wrongWords {
			fmt.Printf("\t%s\n", word)
		}
	}
}

func readBloomFilter(filePath string) *BloomFilter {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	metadata := make([]byte, 16)

	bytesRead, err := reader.Read(metadata)
	if err != nil {
		panic(err)
	} else if bytesRead < 16 {
		panic("invalid file")
	}

	identifer := string(metadata[:4])
	if identifer != IDENTIFIER {
		panic("incorred file identifier")
	}
	version := binary.BigEndian.Uint16(metadata[4:6])
	if version != VERSION_NUMBER {
		panic("incorred version number")
	}
	numHashApplied := binary.BigEndian.Uint16(metadata[6:8])
	totalBits := binary.BigEndian.Uint64(metadata[8:16])

	totalBytes := totalBits / 8

	buffer := make([]byte, 4080)
	bitmap := make([]byte, totalBytes)
	currBytesRead := 0

	for {
		bytesRead, err = reader.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		} else if err == io.EOF {
			for i := 0; i < bytesRead; i++ {
				bitmap[currBytesRead+i] = buffer[i]
			}
			currBytesRead += bytesRead
			break
		}
		for i := 0; i < bytesRead; i++ {
			bitmap[currBytesRead+i] = buffer[i]
		}
		currBytesRead += bytesRead
	}

	if currBytesRead < int(totalBytes) {
		panic("bloom filter missing information")
	}

	return &BloomFilter{
		numHashApplied: uint8(numHashApplied),
		totalBits:      totalBits,
		bloomFilter:    bitmap,
	}
}
