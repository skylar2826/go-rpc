package client

import (
	"context"
	"encoding/binary"
	"net"
)

var DataLength = 8

type connMsg struct {
}

func (c *connMsg) readMsg(ctx context.Context, conn net.Conn) ([]byte, error) {
	dataLengthBs := make([]byte, DataLength)
	_, err := conn.Read(dataLengthBs)
	if err != nil {
		return nil, err
	}
	dataLength := binary.BigEndian.Uint64(dataLengthBs)
	respBs := make([]byte, dataLength)
	_, err = conn.Read(respBs)
	if err != nil {
		return nil, err
	}

	return respBs, nil
}

func (c *connMsg) writeMsg(ctx context.Context, conn net.Conn, data []byte) error {
	reqBs := make([]byte, DataLength+len(data))
	binary.BigEndian.PutUint64(reqBs[:DataLength], uint64(len(data)))
	copy(reqBs[DataLength:], data)

	_, err := conn.Write(reqBs)
	return err
}
