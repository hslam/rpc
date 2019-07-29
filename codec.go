package rpc

import (
	"github.com/golang/protobuf/proto"
	"encoding/json"
	"encoding/xml"
	"encoding/gob"
	"bytes"
)

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) (error)
}

type JsonCodec struct{
}

func (c JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c JsonCodec) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}


type ProtoCodec struct{
}

func (c ProtoCodec) Encode(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (c ProtoCodec) Decode(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}


type XmlCodec struct{
}

func (c XmlCodec) Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (c XmlCodec) Decode(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}


type GobCodec struct{
}

func (c GobCodec) Encode(v interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err :=  gob.NewEncoder(&buffer).Encode(v)
	if err!=nil{
		return nil,err
	}
	return buffer.Bytes(),nil
}

func (c GobCodec) Decode(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(v)
}
