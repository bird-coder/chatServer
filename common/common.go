package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func DealHeader(buf []byte) ([]uint32, error) {
	length := len(buf)
	headers := make([]uint32, 4, 4)
	if length < 16 {
		err := fmt.Errorf("Read msg size failed")
		return headers, err
	}
	pkgBuf := bytes.NewBuffer(buf[0:4])
	pBuf := bytes.NewBuffer(buf[4:8])
	userBuf := bytes.NewBuffer(buf[8:12])
	serverBuf := bytes.NewBuffer(buf[12:16])
	binary.Read(pkgBuf, binary.BigEndian, &headers[0])
	binary.Read(pBuf, binary.BigEndian, &headers[1])
	binary.Read(userBuf, binary.BigEndian, &headers[2])
	binary.Read(serverBuf, binary.BigEndian, &headers[3])
	return headers, nil
}

func createHeader(msg []byte) ([]byte, error) {

}
