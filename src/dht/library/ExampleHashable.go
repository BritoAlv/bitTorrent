package library

type ExampleHashable struct {
	data [NumberBits]uint8
}

func (e ExampleHashable) Hash() [NumberBits]uint8 {
	return e.data
}

func NewExampleHashable(number int) *ExampleHashable {
	return &ExampleHashable{data: IntToBinaryArray(number)}
}
