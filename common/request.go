package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var separator byte = ';'
var metaSeparator byte = ':'

type Request struct {
	/*
		数据长度
		版本号
		序列化协议
		压缩算法
		消息ID

		元数据

		服务名
		方法名
		协议体
	*/
	HeadLength  uint32            `json:"headLength"`
	BodyLength  uint32            `json:"bodyLength"` // 协议版本
	MessageId   uint32            `json:"messageId"`  // 消息ID
	Version     uint8             `json:"version"`
	Compressor  uint8             `json:"compressor"`
	Serializer  uint8             `json:"serializer"`
	Meta        map[string]string `json:"meta"` // 变长头部
	ServiceName string            `json:"serviceName"`
	MethodName  string            `json:"methodName"`
	Data        []byte            `json:"args"`
}

func (r *Request) calculateHeaderLength() {
	// 定长部分 长度为 15
	length := 15
	for key, value := range r.Meta {
		// 2 是分隔符的长度
		// key:value;key2:value2;

		length += len(key) + len(value) + 2
	}

	// 2 是分隔符的长度
	// ... serviceName;methodName;metaKey:metaValue;metaKey2:metaValue2
	r.HeadLength = uint32(length+len(r.ServiceName)+len(r.MethodName)) + 2
}

func (r *Request) Encode() ([]byte, error) {
	r.calculateHeaderLength()
	length := int(r.HeadLength) + len(r.Data)

	data := make([]byte, length)
	binary.BigEndian.PutUint32(data[:4], r.HeadLength)
	binary.BigEndian.PutUint32(data[4:8], r.BodyLength)
	binary.BigEndian.PutUint32(data[8:12], r.MessageId)
	data[12] = r.Version
	data[13] = r.Compressor
	data[14] = r.Serializer

	i := 15 + len(r.ServiceName)
	copy(data[15:i], r.ServiceName)
	data[i] = separator
	j := i + len(r.MethodName) + 1
	copy(data[i+1:j], r.MethodName)
	data[j] = separator

	// 变长部分放最后
	sliceData := data[j+1:]
	for key, value := range r.Meta {
		i = len(key)
		copy(sliceData[0:i], key)

		sliceData[i] = metaSeparator

		j = i + len(value) + 1

		copy(sliceData[i+1:j], value)

		sliceData[j] = separator

		sliceData = sliceData[j+1:]
	}

	copy(data[r.HeadLength:], r.Data)
	return data, nil
}

func (r *Request) Decode(bs []byte) error {
	r.HeadLength = binary.BigEndian.Uint32(bs[:4])
	r.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	r.MessageId = binary.BigEndian.Uint32(bs[8:12])
	r.Version = bs[12]
	r.Compressor = bs[13]
	r.Serializer = bs[14]
	sliceBs := bs[15:]
	idx := bytes.IndexByte(sliceBs, separator)
	if idx > 0 {
		r.ServiceName = string(sliceBs[0:idx])
	}
	sliceBs = sliceBs[idx+1:]
	jdx := bytes.IndexByte(sliceBs, separator)
	if jdx > 0 {
		r.MethodName = string(sliceBs[0:jdx])
	}
	sliceBs = sliceBs[jdx+1:]
	idx = bytes.IndexByte(sliceBs, separator)
	meta := make(map[string]string, 16)
	for {
		if idx == -1 {
			break
		}

		jdx = bytes.IndexByte(sliceBs[0:idx], metaSeparator)
		if jdx == -1 {
			return fmt.Errorf("request header解码错误： meta解析出错 \n")
		}
		key := string(sliceBs[0:jdx])
		value := string(sliceBs[jdx+1 : idx])
		meta[key] = value

		sliceBs = sliceBs[idx+1:]
		idx = bytes.IndexByte(sliceBs, separator)
	}
	r.Meta = meta
	r.Data = bs[r.HeadLength:]
	return nil
}
