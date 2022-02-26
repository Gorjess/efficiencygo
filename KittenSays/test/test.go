package main

import (
	"encoding/json"
	"fmt"
	"github.com/Gorjess/kitten/lib/codec"
	"github.com/Gorjess/kitten/lib/plainstr"
	"strconv"
	"strings"
)

func converter(s string) {
	nums := strings.Split(s, " ")
	bs := make([]byte, len(nums)+1)
	bs[0] = 0

	for i := 1; i < len(nums)+1; i++ {
		b, _ := strconv.Atoi(nums[i-1])
		bs[i] = byte(b)
	}

	obj := plainstr.NewEmpty()
	cc := codec.New()
	er := cc.Unmarshal(bs, obj)
	fmt.Println(er, obj)
}

func jsonTest(s string) {
	// plain json encoder
	ps := plainstr.New(s)
	bs, _ := json.Marshal(ps)
	fmt.Println(bs)

	psEmpty := plainstr.NewEmpty()
	er := json.Unmarshal(bs, psEmpty)
	fmt.Println(er, psEmpty)

	// using codec
	cc := codec.New()
	cbs, _ := cc.Marshal(ps)
	fmt.Println(cbs)

	ccPSE := plainstr.NewEmpty()
	er = cc.Unmarshal(cbs, ccPSE)
	fmt.Println(er, ccPSE)
}

func main() {
	converter("123 34 115 34 58 34 72 101 108 108 111 95 48 34 125")
	//jsonTest("Hello_0")
}
