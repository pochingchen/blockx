package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]byte

// Bytes 地址转字节数组
func (a Address) Bytes() []byte {
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}

	return b
}

// String 地址转十六进制字符串
func (a Address) String() string {
	return hex.EncodeToString(a.Bytes())
}

// AddressFromBytes 由字节数组转地址
func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [20]byte
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}

	return value
}
