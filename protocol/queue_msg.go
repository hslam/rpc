package protocol

import (
	"sync"
)

type Notice struct {
	Id uint16
	RecvChan chan bool
	Recvmsg *Message
	CbChan chan []byte
}

type QueueMsg struct {
	Mu sync.Mutex
	M map[uint16]*Notice
	Queue []uint16
	Pop chan *Notice
}

func(q *QueueMsg)Push(notice *Notice)bool{
	q.Mu.Lock()
	defer q.Mu.Unlock()
	if _,ok:=q.M[notice.Id];ok{
		return false
	}
	q.M[notice.Id]=notice
	q.Queue=append(q.Queue, notice.Id)
	return true
}

func(q *QueueMsg)IsExisted(id uint16)(*Notice, bool){
	q.Mu.Lock()
	defer q.Mu.Unlock()
	if _,ok:=q.M[id];!ok{
		return nil,false
	}
	return q.M[id], true
}

func(q *QueueMsg)SetValue(id uint16,data *Message) {
	q.Mu.Lock()
	defer q.Mu.Unlock()
	if _,ok:=q.M[id];!ok{
		return
	}
	q.M[id].Recvmsg=data
	if len(q.Queue)>0{
		pop_id:=q.Queue[0]
		if _,ok:=q.M[pop_id];ok{
			if q.M[pop_id].Recvmsg!=nil{
				min_msg:=q.M[pop_id]
				delete(q.M,pop_id)
				q.Queue=q.Queue[1:]
				q.Pop<-min_msg
				return
			}
		}
	}
}

func(q *QueueMsg)Check(){
	q.Mu.Lock()
	defer q.Mu.Unlock()
	if len(q.Queue)>0{
		pop_id:=q.Queue[0]
		if _,ok:=q.M[pop_id];ok{
			if q.M[pop_id].Recvmsg!=nil{
				min_msg:=q.M[pop_id]
				delete(q.M,pop_id)
				q.Queue=q.Queue[1:]
				q.Pop<-min_msg
			}
		}
	}
}