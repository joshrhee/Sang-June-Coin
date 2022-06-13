package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func ErrHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	ErrHandler(encoder.Encode(i))
	return aBuffer.Bytes()
}
