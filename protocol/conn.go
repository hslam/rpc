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
		var data []byte
		n, err := reader.Read(buffer)
		if err != nil&& err!=io.ErrShortBuffer {
			goto finish
		}else if err==io.ErrShortBuffer{
			length:=len(buffer)
			if n<=length{
				for{
					length*=2
					tmp_buffer := make([]byte, length)
					n, err=reader.Read(tmp_buffer)
					if err != nil&& err!=io.ErrShortBuffer{
						goto finish
					}else if err==io.ErrShortBuffer{
						continue
					}else {
						data =make([]byte,n)
						copy(data,tmp_buffer[:n])
						tmp_buffer=nil
						break
					}
				}
			}else {
				data = make([]byte, n)
				n, err=reader.Read(data)
				if err!=nil{
					goto finish
				}
			}

		}else {
			data =make([]byte,n)
			copy(data,buffer[:n])
		}
		readChan<-data
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
	func() {
		defer func() {if err := recover(); err != nil {}}()
		finishChan<-true
	}()
endfor:

}
