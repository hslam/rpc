package protocol

import (
	"io"
	"errors"
)

const (
	StreamHeader			= "r"
	StreamHeaderLength		= 1
	StreamBodyLength		= 4
)

func PacketStream(message []byte) []byte {
	var buffer=make([]byte,StreamHeaderLength+StreamBodyLength+len(message))
	copy(buffer[:], []byte(StreamHeader))
	copy(buffer[StreamHeaderLength:], uint32ToBytes(uint32(len(message))))
	copy(buffer[StreamHeaderLength+StreamBodyLength:],message)
	return buffer
	//return append(append([]byte(StreamHeader), uint32ToBytes(uint32(len(message)))...), message...)
}

func UnpackStream(buffer []byte, readerChannel chan []byte) ([]byte,error) {
	length := uint32(len(buffer))
	var i uint32=0
	for i < length{
		if length < i+StreamHeaderLength+StreamBodyLength {
			break
		}
		if string(buffer[i:i+StreamHeaderLength]) == StreamHeader {
			bodyLength := bytesToUint32(buffer[i+StreamHeaderLength: i+StreamHeaderLength+StreamBodyLength])
			if length < i+StreamHeaderLength+StreamBodyLength+bodyLength {
				break
			}
			data := buffer[i+StreamHeaderLength+StreamBodyLength : i+StreamHeaderLength+StreamBodyLength+bodyLength]
			readerChannel <- data
			i += StreamHeaderLength+StreamBodyLength + bodyLength
		}else {
			return nil,errors.New("invalid")
		}
	}
	if i == length {
		return make([]byte, 0),nil
	}
	return buffer[i:],nil
}

func ReadStream(reader io.Reader, readChan chan []byte, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	tmpBuffer := make([]byte, 0)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			goto finish
		}
		if n>0{
			var err error
			tmpBuffer,err = UnpackStream(append(tmpBuffer, buffer[:n]...), readChan)
			if err != nil {
				goto finish
			}
		}
		select {
		case <-stopChan:
			goto endfor
		default:
		}
	}
finish:
	finishChan<-true
endfor:
}

func WriteStream(writer io.Writer, writeChan chan []byte, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for send_data:= range writeChan{
		dataPacket:=PacketStream(send_data)
		_, err := writer.Write(dataPacket)
		if err != nil {
			goto finish
		}
		select {
		case <-stopChan:
			goto endfor
		default:
		}
	}
finish:
	finishChan<-true
endfor:
}
