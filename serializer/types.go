package serializer

type Serializer interface {
	Code() uint8
	Encode(data interface{}) ([]byte, error)
	DeCode(bs []byte, data interface{}) error
}

var _ Serializer = &Json{}
