package common

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

func BitOn(byteField []byte, bitIdx int) bool {

	byteIdx := bitIdx / 8
	bitIdxOnCurrByte := bitIdx % 8

	if len(byteField) <= byteIdx {
		return false
	}
	return byteField[byteIdx]&(1<<bitIdxOnCurrByte) != 0
}

func TurnBitOn(byteField []byte, bitIdx int) {

	byteIdx := bitIdx / 8
	bitIdxOnCurrByte := bitIdx % 8

	byteField[byteIdx] |= 1 << bitIdxOnCurrByte
}

func BitHashChecker(buf []byte, pieceHash [20]byte) bool {
	hash := sha1.Sum(buf)
	if bytes.Equal(hash[:], pieceHash[:]) {
		return true
	}
	fmt.Println("Hash mismatch")
	return false
}

func ConvertStringToByteArray(str string) *[20]byte {
	var bytes [20]byte
	copy(bytes[:], []byte(str))
	return &bytes
}

func CalcMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}