package library

import (
	"bittorrent/common"
	"math/rand/v2"
	"sort"
	"strconv"
	"time"
)

var usedId = map[ChordHash]bool{}

func generateTaskId() int64 {
	return rand.Int64()
}

func generateRandomKey() ChordHash {
	return rand.Int64() % (1 << NumberBits)
}

func Sort(ids []ChordHash) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
}

func GenerateRandomBinaryId() ChordHash {
	var result ChordHash
	for {
		result = generateRandomKey()
		if _, exist := usedId[result]; !exist {
			usedId[result] = true
			break
		}
	}
	return result
}

// Between : starting from L + 1 in a clockwise order, I can reach M before R + 1.
func Between(L ChordHash, M ChordHash, R ChordHash) bool {
	L = (L + 1) % (1 << NumberBits)
	for {
		if L == M {
			return true
		}
		if L == R {
			break
		}
		L = (L + 1) % (1 << NumberBits)
	}
	return false
}

func SetLogDirectoryPath(name string) {
	common.LogsPath = "./logs/" + name + strconv.Itoa(time.Now().Nanosecond()) + "/"
}
