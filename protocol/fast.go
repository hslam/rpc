package protocol

import (
	"github.com/valyala/fasthttp"
	"net/url"
	"math/rand"
)

func HandleFASTHTTP(fastclient *fasthttp.Client,address string,readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	WindowSize:=64
	readMessageChan := make(chan *Message,WindowSize)
	idChan := make(chan uint16,WindowSize)
	queueMsg:=&QueueMsg{
		M:make(map[uint16]*Notice),
		Queue:make([]uint16,0),
		Pop:make(chan *Notice,WindowSize),
	}
	var startbit =uint(rand.Intn(13))
	var id =uint16(rand.Int31n(int32(1<<startbit)))
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
			case old_notice,ok := <-queueMsg.Pop:
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
		id++
		if id>1<<15{
			id=1
		}
		idChan<-id
		notice:=&Notice{Id:id,}
		queueMsg.Push(notice)
		msg:=&Message{OprationTypeData,id,send_data}
		go func(fastclient *fasthttp.Client,address string,msg *Message) {
			for {
				req := &fasthttp.Request{}
				req.Header.SetMethod("POST")
				req.SetBody(msg.message)
				u := url.URL{Scheme: "http", Host: address, Path: "/"}
				req.SetRequestURI(u.String())
				resp := &fasthttp.Response{}
				err := fastclient.Do(req, resp)
				if err != nil {
					continue
				}
				if len(resp.Body())>0{
					readMessageChan<-&Message{OprationTypeData,msg.id,resp.Body()}
				}else {
					readMessageChan<-&Message{OprationTypeAck,msg.id,nil}
				}
				goto endfor
			}
			stopChan<-true
			endfor:
		}(fastclient,address,msg)
	}
	close(readMessageChan)
	close(idChan)
	close(queueMsg.Pop)
	queueMsg.Queue=nil
	queueMsg.M=nil
	queueMsg=nil
}

