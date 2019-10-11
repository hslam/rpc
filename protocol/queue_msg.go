package protocol

import (
	"sync"
	"fmt"
)

type Notice struct {
	Id uint16
	RecvChan chan bool
	Recvmsg *Message
	CbChan chan []byte
}

type QueueMsg struct {
	mu sync.Mutex
	m map[uint16]*Notice
	queue []uint16
	front int
	rear  int
	maxSize int
	pop chan *Notice

}
func NewQueueMsg(size int) *QueueMsg {
	return &QueueMsg{
		m:make(map[uint16]*Notice),
		queue:make([]uint16,size),
		pop:make(chan *Notice,size),
		front: 0,
		rear:  0,
		maxSize:size,
	}
}
func(q *QueueMsg)Push(notice *Notice)bool{
	q.mu.Lock()
	defer q.mu.Unlock()
	if _,ok:=q.m[notice.Id];ok{
		return false
	}
	q.m[notice.Id]=notice
	q.Enqueue(notice.Id)
	return true
}
func(q *QueueMsg)Pop()chan *Notice{
	return q.pop
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
	q.m[id].Recvmsg=data
	if q.Length()>0{
		pop_id,err:=q.Front()
		if err==nil{
			if _,ok:=q.m[pop_id];ok{
				if q.m[pop_id].Recvmsg!=nil{
					min_msg:=q.m[pop_id]
					delete(q.m,pop_id)
					q.Dequeue()
					q.pop<-min_msg
					return
				}
			}
		}
	}
}

func(q *QueueMsg)Check(){
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Length()>0{
		pop_id,err:=q.Front()
		if err==nil{
			if _,ok:=q.m[pop_id];ok{
				if q.m[pop_id].Recvmsg!=nil{
					min_msg:=q.m[pop_id]
					delete(q.m,pop_id)
					q.Dequeue()
					q.pop<-min_msg
					return
				}
			}
		}
	}
}


func (q *QueueMsg) Length() int {
	return (q.rear - q.front + q.maxSize) % q.maxSize
}

func (q *QueueMsg) Enqueue(e uint16) error {
	if (q.rear+1)%q.maxSize == q.front {
		return fmt.Errorf("quque is full")
	}
	q.queue[q.rear] = e
	q.rear = (q.rear + 1) % q.maxSize
	return nil
}

func (q *QueueMsg) Dequeue() (e uint16, err error) {
	if q.rear == q.front {
		return e, fmt.Errorf("quque is empty")
	}
	e = q.queue[q.front]
	q.front = (q.front + 1) % q.maxSize
	return e, nil
}

func (q *QueueMsg) Front() (e uint16, err error) {
	if q.rear == q.front {
		return e, fmt.Errorf("quque is empty")
	}
	e = q.queue[q.front]
	return e, nil
}

func (q *QueueMsg)Close(){
	close(q.pop)
	q.queue=nil
	q.m=nil
	q=nil
}