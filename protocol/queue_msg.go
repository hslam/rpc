package protocol

import (
	"sync"
)

type Notice struct {
	id uint16
	recvChan chan bool
	recvmsg *Message
}

type QueueMsg struct {
	mu sync.Mutex
	m map[uint16]*Notice
	queue []uint16
	pop chan *Notice
}

func(q *QueueMsg)Push(notice *Notice)bool{
	q.mu.Lock()
	defer q.mu.Unlock()
	if _,ok:=q.m[notice.id];ok{
		return false
	}
	q.m[notice.id]=notice
	q.queue=append(q.queue, notice.id)
	return true
}

func(q *QueueMsg)IsExisted(id uint16)(*Notice, bool){
	q.mu.Lock()
	defer q.mu.Unlock()
	if _,ok:=q.m[id];!ok{
		return nil,false
	}
	return q.m[id], true
}

func(q *QueueMsg)SetValue(id uint16,data *Message) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _,ok:=q.m[id];!ok{
		return
	}
	q.m[id].recvmsg=data
	if len(q.queue)>0{
		pop_id:=q.queue[0]
		if _,ok:=q.m[pop_id];ok{
			if q.m[pop_id].recvmsg!=nil{
				min_msg:=q.m[pop_id]
				delete(q.m,pop_id)
				q.queue=q.queue[1:]
				q.pop<-min_msg
				return
			}
		}
	}
}

func(q *QueueMsg)Check(){
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue)>0{
		pop_id:=q.queue[0]
		if _,ok:=q.m[pop_id];ok{
			if q.m[pop_id].recvmsg!=nil{
				min_msg:=q.m[pop_id]
				delete(q.m,pop_id)
				q.queue=q.queue[1:]
				q.pop<-min_msg
			}
		}
	}
}