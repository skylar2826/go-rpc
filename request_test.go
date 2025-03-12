package main

import (
	"fmt"
	"go-geektime3/common"
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	r := common.Request{
		BodyLength: uint32(1),
		MessageId:  uint32(2),
		Version:    uint8(3),
		Compressor: uint8(4),
		Serializer: uint8(5),
		Meta: map[string]string{
			"meta1": "value1",
			"meta2": "value2",
		},
		ServiceName: "serviceName",
		MethodName:  "methodName",
		Data:        []byte("hello world"),
	}
	bs, err := r.Encode()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(bs)

	r1 := &common.Request{}
	err = r1.Decode(bs)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(r1)
}
