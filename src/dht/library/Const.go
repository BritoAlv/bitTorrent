package library

const NumberBits = 8
const WaitingTime = 1
const Attempts = 5

type ChordHash = [NumberBits]uint8
type Store = map[ChordHash][]byte
