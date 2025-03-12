package main

import (
	"fmt"
	"go-geektime3/common"
	"log"
	"testing"
)

func TestResponse(t *testing.T) {
	r := common.Response{
		BodyLength: uint32(1),
		MessageId:  uint32(2),
		Version:    uint8(3),
		Compressor: uint8(4),
		Serializer: uint8(5),
		Error:      []byte("error"),
		Data:       []byte("hello world"),
	}
	bs, err := r.Encode()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(bs)

	r1 := &common.Response{}
	err = r1.Decode(bs)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(r1)
}
