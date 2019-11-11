package protocol

import (
	"io"
	"errors"
	"time"
	"sync"
)

const (
	StreamHeader			= "r"
	StreamHeaderLength		= 1
	StreamBodyLength		= 4
)

func PacketStream(body []byte) []byte {
	var buffer=make([]byte,StreamHeaderLength+StreamBodyLength+len(body))
	copy(buffer[:], []byte(StreamHeader))
	copy(buffer[StreamHeaderLength:], uint32ToBytes(uint32(len(body))))
	copy(buffer[StreamHeaderLength+StreamBodyLength:],body)
	return buffer
	//return append(append([]byte(StreamHeader), uint32ToBytes(uint32(len(body)))...), body...)
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
			body := buffer[i+StreamHeaderLength+StreamBodyLength : i+StreamHeaderLength+StreamBodyLength+bodyLength]
			readerChannel <- body
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
	func() {
		defer func() {if err := recover(); err != nil {}}()
		finishChan<-true
	}()
endfor:
}

func WriteStream(writer io.Writer, writeChan chan []byte, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	useBuffer:=false
	var ticker *time.Ticker
	var ch chan bool
	if useBuffer{
		ticker =time.NewTicker(time.Microsecond)
		buf:=make([]byte,0)
		var mu sync.Mutex
		ch=make(chan bool,1)
		go func() {
			for send_data:= range writeChan{
				dataPacket:=PacketStream(send_data)
				func(){
					mu.Lock()
					defer mu.Unlock()
					buf=append(buf,dataPacket...)
					if len(buf)>1024{
						writer.Write(buf)
						buf=make([]byte,0)
					}
				}()
			}
		}()
		for range ticker.C{
			select {
			case <-ch:
				goto endfor
			default:
			}
			err:=func()error{
				mu.Lock()
				defer mu.Unlock()
				if len(buf)>0{
					_, err := writer.Write(buf)
					buf=make([]byte,0)
					return err
				}
				return nil
			}()
			if err != nil {
				goto finish
			}
			select {
			case <-stopChan:
				goto endfor
			default:
			}
		}
	}else {
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
	}
finish:
	func() {
		defer func() {if err := recover(); err != nil {}}()
		finishChan<-true
	}()
endfor:
	if useBuffer{
		ticker.Stop()
		ticker=nil
		ch<-true
	}
}
