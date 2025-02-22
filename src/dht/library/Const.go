package library

import "strconv"

const NumberBits = 8
const WaitingTime = 1
const Attempts = 5
const StateQueryWaitTime = 1

type ChordHash = int64
type Store = map[ChordHash][]byte

func ToString(A ChordHash) string {
	return strconv.Itoa(int(A))
}
