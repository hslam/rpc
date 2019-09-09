package protocol

import (
	"io"
	"errors"
	"fmt"
	"time"
	"math/rand"
	"hslam.com/mgit/Mort/rpc/log"
)

type  OprationType byte

const (
	MessageHeaderOprationLength   = 1
	MessageHeaderIDLength   = 2
	MessageHeaderCheckSumLength   = 2
	MessageHeaderLength   = 5

	OprationTypeInvalid  OprationType= 0
	OprationTypeData  OprationType= 1
	OprationTypeAck  OprationType= 2
	OprationTypeHeartbeat  OprationType= 3
)

func PacketMessage(oprationType OprationType,id uint16,message []byte) []byte {
	switch oprationType {
	case OprationTypeData:
		header:=make([]byte,MessageHeaderLength)
		header[0]=byte(OprationTypeData)
		idBytes:=uint16ToBytes(id)
		header[1],header[2]=idBytes[0],idBytes[1]
		header[3],header[4]=byte(0),byte(0)
		msg:=append(header, message...)
		checkSumBuffer:=make([]byte,len(msg))
		copy(checkSumBuffer,msg)
		checkSumBytes:=uint16ToBytes(checkSum(checkSumBuffer))
		msg[3],msg[4]=checkSumBytes[0],checkSumBytes[1]
		return msg
	case OprationTypeAck:
		header:=make([]byte,MessageHeaderLength)
		header[0]=byte(OprationTypeAck)
		idBytes:=uint16ToBytes(id)
		header[1],header[2]=idBytes[0],idBytes[1]
		header[3],header[4]=byte(0),byte(0)
		checkSumBuffer:=make([]byte,len(header))
		copy(checkSumBuffer,header)
		checkSumBytes:=uint16ToBytes(checkSum(checkSumBuffer))
		header[3],header[4]=checkSumBytes[0],checkSumBytes[1]
		return header
	case OprationTypeHeartbeat:
		header:=make([]byte,MessageHeaderLength)
		header[0]=byte(OprationTypeHeartbeat)
		idBytes:=uint16ToBytes(id)
		header[1],header[2]=idBytes[0],idBytes[1]
		header[3],header[4]=byte(0),byte(0)
		checkSumBuffer:=make([]byte,len(header))
		copy(checkSumBuffer,header)
		checkSumBytes:=uint16ToBytes(checkSum(checkSumBuffer))
		header[3],header[4]=checkSumBytes[0],checkSumBytes[1]
		return header
	}
	return nil
}

func UnpackMessage(buffer []byte)(OprationType,uint16,[]byte,error){
	oprationType:=OprationType(buffer[0])
	switch oprationType {
	case OprationTypeData:
		idBytes:=buffer[MessageHeaderOprationLength:MessageHeaderOprationLength+MessageHeaderIDLength]
		checkSumBytes:=buffer[MessageHeaderOprationLength+MessageHeaderIDLength:MessageHeaderLength]
		checkSumValue:=bytesToUint16(checkSumBytes)
		checkSumBuffer:=make([]byte,len(buffer))
		copy(checkSumBuffer,buffer)
		checkSumBuffer[3],checkSumBuffer[4]=byte(0),byte(0)
		if checkSum(checkSumBuffer) != checkSumValue{
			return OprationTypeInvalid, 0,nil,errors.New("message data is error"+fmt.Sprintln( checkSum(checkSumBuffer) ,checkSumValue,checkSumBytes))
		}
		return OprationTypeData, bytesToUint16(idBytes),buffer[MessageHeaderLength:],nil
	case OprationTypeAck:
		idBytes:=buffer[MessageHeaderOprationLength:MessageHeaderOprationLength+MessageHeaderIDLength]
		checkSumBytes:=buffer[MessageHeaderOprationLength+MessageHeaderIDLength:MessageHeaderLength]
		checkSumValue:=bytesToUint16(checkSumBytes)
		checkSumBuffer:=make([]byte,len(buffer))
		copy(checkSumBuffer,buffer)
		checkSumBuffer[3],checkSumBuffer[4]=byte(0),byte(0)
		if checkSum(checkSumBuffer) != checkSumValue{
			return OprationTypeInvalid, 0,nil,errors.New("message ack is error"+fmt.Sprintln( checkSum(checkSumBuffer) ,checkSumValue,checkSumBytes))
		}
		return OprationTypeAck, bytesToUint16(idBytes),nil,nil
	case OprationTypeHeartbeat:
		idBytes:=buffer[MessageHeaderOprationLength:MessageHeaderOprationLength+MessageHeaderIDLength]
		checkSumBytes:=buffer[MessageHeaderOprationLength+MessageHeaderIDLength:MessageHeaderLength]
		checkSumValue:=bytesToUint16(checkSumBytes)
		checkSumBuffer:=make([]byte,len(buffer))
		copy(checkSumBuffer,buffer)
		checkSumBuffer[3],checkSumBuffer[4]=byte(0),byte(0)
		if checkSum(checkSumBuffer) != checkSumValue{
			return OprationTypeInvalid, 0,nil,errors.New("message heartbeat is error"+fmt.Sprintln( checkSum(checkSumBuffer) ,checkSumValue,checkSumBytes))
		}
		return OprationTypeHeartbeat, bytesToUint16(idBytes),nil,nil
	default:
		return OprationTypeInvalid,0,nil,errors.New("message is error")
	}
}

type Message struct {
	oprationType OprationType
	id uint16
	message []byte
}

func ReadMessage(reader io.Reader, readMessageChan chan *Message, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buffer := make([]byte, 1024*64)
	for {

		n, err := reader.Read(buffer[:])
		if err != nil {
			goto finish
		}
		var data =make([]byte,n)
		copy(data,buffer[:n])
		log.Debugf("recv data size %d",n)
		oprationType,id,msg,err:=UnpackMessage(data)
		if err != nil {
			continue
		}
		if oprationType==OprationTypeAck|| oprationType==OprationTypeData{
			msg:=&Message{oprationType,id,msg}
			readMessageChan<-msg
			continue
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

func WriteMessage(writer io.Writer, writeMessageChan chan *Message, stopChan chan bool,finishChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for msg:= range writeMessageChan{
		data:=PacketMessage(msg.oprationType,msg.id,msg.message)
		log.Debugf("send data size %d",len(data))
		_, err := writer.Write(data)
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

func HandleMessage(readWriter io.ReadWriter,readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	WindowSize:=1024
	idChan := make(chan uint16,WindowSize)
	queueMsg:=NewQueueMsg(WindowSize)
	var startbit =uint(rand.Intn(13))
	var id =uint16(rand.Int31n(int32(1<<startbit)))
	var max_id=uint16(1<<16-1)
	var SmoothedRTT int64=1000000
	rto := &RTO{}
	lastRTO:=rto.updateRTT(SmoothedRTT)
	readMessageChan := make(chan *Message,WindowSize)
	writeMessageChan := make(chan *Message,WindowSize)
	finishMessageChan:=make(chan bool)
	stopReadMessageChan := make(chan bool,1)
	stopWriteMessageChan := make(chan bool,1)
	go ReadMessage(readWriter, readMessageChan, stopReadMessageChan,finishMessageChan)
	go WriteMessage(readWriter,writeMessageChan, stopWriteMessageChan,finishMessageChan)
	go func(readMessageChan chan *Message, queueMsg *QueueMsg) {
		for{
			select {
			case revmsg,ok := <-readMessageChan:
				if ok{
					if notice,ok:= queueMsg.IsExisted(revmsg.id);ok{
						notice.RecvChan<-true
						queueMsg.SetValue(revmsg.id,revmsg)
					}
				}else {
					goto endfor
				}
			}
		}
	endfor:
	}(readMessageChan,queueMsg)
	go func(idChan chan uint16,queueMsg *QueueMsg,readChan chan []byte) {
		for{
			select {
			case old_notice,ok := <-queueMsg.Pop():
				if ok{
					if old_notice.Recvmsg.oprationType==OprationTypeData{
						readChan<-old_notice.Recvmsg.message
						<-idChan
					}else if old_notice.Recvmsg.oprationType==OprationTypeAck{
						<-idChan
					}
					queueMsg.Check()
				}else {
					goto endfor
				}
			}
		}
	endfor:
	}(idChan,queueMsg,readChan)
	for send_data:= range writeChan{
		id=(id+1)%max_id
		idChan<-id
		recvChan:=make(chan bool)
		notice:=&Notice{Id:id,RecvChan:recvChan}
		queueMsg.Push(notice)
		msg:=&Message{OprationTypeData,id,send_data}
		go func(msg *Message,rto *RTO,recvChan chan bool ) {
			start_time:=time.Now().UnixNano()
			timer := time.NewTimer(time.Microsecond*time.Duration(lastRTO*8))
			for {
				func(){
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					writeMessageChan <-msg
				}()
				select {
				case <-recvChan:
					goto endfor
				case <-timer.C:
				}
			}
		endfor:
			RTT:=(time.Now().UnixNano()-start_time)/1000
			lastRTO = rto.updateRTT(RTT)
		}(msg,rto,recvChan)
		select {
		case <-stopChan:
			stopReadMessageChan<-true
			stopWriteMessageChan<-true
			goto endfor
		case <-finishMessageChan:
			finishChan<-true
			stopReadMessageChan<-true
			stopWriteMessageChan<-true
			goto endfor
		default:
		}
	}
endfor:
	close(readMessageChan)
	close(writeMessageChan)
	close(finishMessageChan)
	close(stopWriteMessageChan)
	close(stopReadMessageChan)
	close(idChan)
	queueMsg.Close()
	rto=nil
	log.Traceln("HandleMessage end")

}