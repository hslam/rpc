package protocol

import (
	"io"
)

const (
	StreamHeaderLength   = 4
)

func PacketStream(message []byte) []byte {
	return append(uint32ToBytes(uint32(len(message))), message...)
}

func UnpackStream(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)
	var i int=0
	for i < length{
		if length < i+StreamHeaderLength {
			break
		}
		messageLength := int(bytesToUint32(buffer[i : i+StreamHeaderLength]))
		if length < i+StreamHeaderLength+messageLength {
			break
		}
		data := buffer[i+StreamHeaderLength : i+StreamHeaderLength+messageLength]
		readerChannel <- data
		i += StreamHeaderLength + messageLength
	}
	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
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
			tmpBuffer = UnpackStream(append(tmpBuffer, buffer[:n]...), readChan)
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
