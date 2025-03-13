package common

import (
	"encoding/binary"
)

type Response struct {
	/**
	数据长度
	序列化协议
	压缩算法
	消息ID
	错误信息
	协议体
	*/
	HeadLength uint32 `json:"headLength"`
	BodyLength uint32 `json:"bodyLength"`
	MessageId  uint32 `json:"messageId"`
	Version    uint8  `json:"version"`
	Compressor uint8  `json:"compressor"`
	Serializer uint8  `json:"serializer"`
	Error      []byte `json:"error"`
	Data       []byte `json:"data"`
}

func (r *Response) calculateHeaderLength() {
	// 定长部分 长度为 15
	length := 15
	r.HeadLength = uint32(length + len(r.Error))
}

func (r *Response) Encode() ([]byte, error) {
	r.calculateHeaderLength()
	length := int(r.HeadLength) + len(r.Data)

	data := make([]byte, length)
	binary.BigEndian.PutUint32(data[:4], r.HeadLength)
	binary.BigEndian.PutUint32(data[4:8], r.BodyLength)
	binary.BigEndian.PutUint32(data[8:12], r.MessageId)
	data[12] = r.Version
	data[13] = r.Compressor
	data[14] = r.Serializer

	copy(data[15:r.HeadLength], r.Error)
	copy(data[r.HeadLength:], r.Data)
	return data, nil
}

func (r *Response) Decode(bs []byte) error {
	r.HeadLength = binary.BigEndian.Uint32(bs[:4])
	r.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	r.MessageId = binary.BigEndian.Uint32(bs[8:12])
	r.Version = bs[12]
	r.Compressor = bs[13]
	r.Serializer = bs[14]
	if r.HeadLength > 14 {
		r.Error = bs[15:r.HeadLength]
	}

	r.Data = bs[r.HeadLength:]
	return nil
}
