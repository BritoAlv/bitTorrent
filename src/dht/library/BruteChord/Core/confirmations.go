package Core

type confirmations struct {
	Confirmation bool
	Value        []byte // this have to be interpreted by the caller by now.
}
