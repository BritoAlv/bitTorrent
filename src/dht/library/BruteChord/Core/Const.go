package Core

const NumberBits = 8
const WaitingTime = 1
const Attempts = 5

type ChordHash = int64
type Store = map[ChordHash][]byte
