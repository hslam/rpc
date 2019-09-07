package protocol

import (
	"io"
)

func ReadConn(reader io.Reader, readChan chan []byte, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buffer := make([]byte, 65536)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			goto finish
		}
		var data =make([]byte,n)
		copy(data,buffer[:n])
		readChan<-data
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

func WriteConn(writer io.Writer, writeChan chan []byte, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for send_data:= range writeChan{
		_, err := writer.Write(send_data)
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
