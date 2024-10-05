package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	filePath := "./resources/words.txt"
	var maxSize float64 = 470000
	desiredErrorRate := 0.001

	testRandomInserts(filePath, maxSize, desiredErrorRate)
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
