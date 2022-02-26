package codec

import "encoding/json"

type ICodec interface {
	Marshal(input interface{}) ([]byte, error)
	Unmarshal(input []byte, output interface{}) error
}

type IContent interface {
	String() string
}

type Codec struct {
}

func New() Codec {
	return Codec{}
}

func (c Codec) Marshal(input IContent) ([]byte, error) {
	// do marshal
	bs, er := json.Marshal(input)
	if er != nil {
		return nil, er
	}
	l := len(bs)

	// combine head and body
	buf := make([]byte, 1, 1+l)
	buf[0] = byte(l)
	buf = append(buf, bs...)

	return buf, nil
}

func (c Codec) Unmarshal(input []byte, output interface{}) error {
	return json.Unmarshal(input, output)
}
