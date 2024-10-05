package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
)

func testRandomInserts(filePath string, maxSize float64, desiredErrorRate float64) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	bloomFilter := createBloomFilter(maxSize, desiredErrorRate)

	didInsertWordMap := make(map[string]bool)

	for {
		buffer, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		word := string(buffer)

		randomNumber := rand.Uint32()

		if randomNumber < ^uint32(0)/2 {
			bloomFilter.insert(word)
			didInsertWordMap[word] = true
		} else {
			didInsertWordMap[word] = false
		}
	}

	file.Close()

	falsePositives := 0
	falseNegatives := 0

	for word, didInsert := range didInsertWordMap {
		doesWordExist := bloomFilter.query(word)

		if !didInsert && doesWordExist {
			falsePositives += 1
		} else if didInsert && !doesWordExist {
			falseNegatives += 1
		}
	}

	falsePositiveRate := float32(falsePositives) / float32(len(didInsertWordMap))
	falsePositiveRate *= 100

	fmt.Printf("The false positive rate is %f percent with %d false positives and %d total inserts\n", falsePositiveRate, falsePositives, bloomFilter.insertedElements)
	fmt.Printf("There are %d false negatives\n", falseNegatives)
}
