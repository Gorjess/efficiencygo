package plainstr

type PlainStr struct {
	S string `json:"s"`
}

func New(s string) *PlainStr {
	return &PlainStr{S: s}
}

func NewEmpty() *PlainStr {
	return &PlainStr{}
}

func (ps *PlainStr) String() string {
	return ps.S
}

func (ps *PlainStr) Size() int {
	return len(ps.S)
}
