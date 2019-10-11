package protocol

import (
	"math/rand"
	"fmt"
)

type SyncConn interface {
	Do(reqBody []byte)([]byte,error)
}

func HandleSyncConn(syncConn SyncConn,recvChan chan []byte,sendChan chan []byte, stopChan chan bool,windowSize int){
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	pipelineChan := make(chan bool,windowSize)
	readMessageChan := make(chan *Message,windowSize)
	idChan := make(chan uint16,windowSize)
	queueMsg:=NewQueueMsg(windowSize)
	var startbit =uint(rand.Intn(5))
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
						func(){
							defer func() {
								if err := recover(); err != nil {
									fmt.Println(err)
								}
							}()
							readChan<-old_notice.Recvmsg.message
						}()
						<-idChan
					}else if old_notice.Recvmsg.oprationType==OprationTypeAck{
						<-idChan
					}
					func(){
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						queueMsg.Check()
					}()
				}else {
					goto endfor
				}
			}
		}
	endfor:
	}(idChan,queueMsg,recvChan)
	for send_data:= range sendChan{
		pipelineChan<-true
		id=(id+1)%max_id
		idChan<-id
		notice:=&Notice{Id:id,}
		queueMsg.Push(notice)
		msg:=&Message{OprationTypeData,id,send_data}
		go func(syncConn SyncConn,msg *Message,readMessageChan chan *Message) {
			for {
				respBody,err := syncConn.Do(msg.message)
				if err != nil {
					continue
				}
				if len(respBody)>0{

					func(){
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						readMessageChan<-&Message{OprationTypeData,msg.id,respBody}
					}()
				}else {
					func(){
						defer func() {
							if err := recover(); err != nil {

							}
						}()
						readMessageChan<-&Message{OprationTypeAck,msg.id,nil}
					}()
				}
				goto endfor
			}
		endfor:
			<-pipelineChan
		}(syncConn,msg,readMessageChan)
		select {
		case <-stopChan:
			goto finish
		default:
		}
	}
finish:
	close(readMessageChan)
	close(idChan)
	queueMsg.Close()
}
