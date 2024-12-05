package wsmux

import (
	"bytes"
	"encoding/gob"
)

type WsPackage struct {
	ID      uint16 `json:"id"`
	SeqN    uint16 `json:"seq_n"`
	Message []byte `json:"message"`

	err error
}

func (wp *WsPackage) ToBytes() []byte {
	data, err := ToBytes(wp)
	if err != nil {
		wp.err = err
	}
	return data
}

func (wp *WsPackage) FromBytes(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wp)
	if err != nil {
		wp.err = err
	}
	return err
}

func ToBytes(wp *WsPackage) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(wp)

	return buf.Bytes(), err
}

func FromBytes(data []byte) (*WsPackage, error) {
	var wp WsPackage
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&wp)
	if err != nil {
		return nil, err
	}
	return &wp, nil
}
