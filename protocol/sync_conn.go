package protocol

import (
	"math/rand"
)

type SyncConn interface {
	Do(reqBody []byte)([]byte,error)
}

func HandleSyncConn(syncConn SyncConn,readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	WindowSize:=64
	readMessageChan := make(chan *Message,WindowSize)
	idChan := make(chan uint16,WindowSize)
	queueMsg:=NewQueueMsg(WindowSize)
	var startbit =uint(rand.Intn(13))
	var id =uint16(rand.Int31n(int32(1<<startbit)))
	var max_id=uint16(1<<16-1)
	go func(readMessageChan chan *Message, queueMsg *QueueMsg) {
		for{
			select {
			case revmsg,ok := <-readMessageChan:
				if ok{
					if _,ok:= queueMsg.IsExisted(revmsg.id);ok{
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
		notice:=&Notice{Id:id,}
		queueMsg.Push(notice)
		msg:=&Message{OprationTypeData,id,send_data}
		go func(syncConn SyncConn,msg *Message) {
			for {
				respBody,err := syncConn.Do(msg.message)
				if err != nil {
					continue
				}
				if len(respBody)>0{
					readMessageChan<-&Message{OprationTypeData,msg.id,respBody}
				}else {
					readMessageChan<-&Message{OprationTypeAck,msg.id,nil}
				}
				goto endfor
			}
			stopChan<-true
			endfor:
		}(syncConn,msg)
	}
	close(readMessageChan)
	close(idChan)
	queueMsg.Close()
}
