package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

const FNV_PRIME uint64 = 1099511628211
const FNV_OFFSET_BASIS uint64 = 14695981039346656037

type BloomFilter struct {
	numHashApplied   uint8
	totalBits        uint64
	bloomFilter      []byte
	insertedElements uint64
	insertCollisions uint64
}

func (b *BloomFilter) String() string {
	return fmt.Sprintf(
		`number of applied hash functions: %d
total length of bloom filter: %d
current value of bloom filter: %d
inserted elements: %d
number of insert collisions: %d`,
		b.numHashApplied,
		b.totalBits,
		b.bloomFilter,
		b.insertedElements,
		b.insertCollisions,
	)
}

func (b *BloomFilter) insert(word string) {

	prehash := []byte(word)
	posthash := fnv1_hash(prehash)
	bitPosition := posthash % b.totalBits
	index := bitPosition / 8
	indexBitPosition := byte(bitPosition % 8)

	b.bloomFilter[index] |= (1 << indexBitPosition)

	prehash = make([]byte, 8)
	binary.LittleEndian.PutUint64(prehash, posthash)

	for i := uint8(0); i < b.numHashApplied-1; i++ {
		posthash := fnv1_hash(prehash)
		bitPosition := posthash % b.totalBits
		index := bitPosition / 8
		indexBitPosition := byte(bitPosition % 8)

		b.bloomFilter[index] |= (1 << indexBitPosition)

		binary.LittleEndian.PutUint64(prehash, posthash)
	}
	b.insertedElements += 1
}

func (b *BloomFilter) query(word string) bool {
	prehash := []byte(word)
	posthash := fnv1_hash(prehash)
	bitPosition := posthash % b.totalBits
	index := bitPosition / 8
	indexBitPosition := byte(bitPosition % 8)

	if b.bloomFilter[index]&(1<<indexBitPosition) == 0 {
		return false
	}

	prehash = make([]byte, 8)
	binary.LittleEndian.PutUint64(prehash, posthash)

	for i := uint8(0); i < b.numHashApplied-1; i++ {
		posthash := fnv1_hash(prehash)
		bitPosition := posthash % b.totalBits
		index := bitPosition / 8
		indexBitPosition := byte(bitPosition % 8)

		if b.bloomFilter[index]&(1<<indexBitPosition) == 0 {
			return false
		}

		binary.LittleEndian.PutUint64(prehash, posthash)
	}
	return true
}

func fnv1_hash(word []byte) uint64 {
	var hash uint64 = FNV_OFFSET_BASIS

	for i := 0; i < len(word); i++ {
		hash *= FNV_PRIME
		hash ^= uint64(word[i])
	}

	return hash
}

func createBloomFilter(numElements float64, errorRate float64) *BloomFilter {
	totalBits := -(numElements * math.Log(float64(errorRate))) / math.Pow(math.Log(2), 2)
	numHashApplied := totalBits / numElements * math.Log(2)
	bloomFilterSize := uint32(math.Floor(totalBits / 8))

	return &BloomFilter{
		numHashApplied:   uint8(numHashApplied),
		totalBits:        uint64(totalBits),
		bloomFilter:      make([]byte, bloomFilterSize),
		insertedElements: 0,
		insertCollisions: 0,
	}
}
