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
	Mu sync.Mutex
	M map[uint16]*Notice
	Queue []uint16
	front int
	rear  int
	maxSize int
	Pop chan *Notice

}
func NewQueueMsg(WindowSize int) *QueueMsg {
	return &QueueMsg{
		M:make(map[uint16]*Notice),
		Queue:make([]uint16,WindowSize),
		Pop:make(chan *Notice,WindowSize),
		front: 0,
		rear:  0,
		maxSize:WindowSize,
	}
}
func(q *QueueMsg)Push(notice *Notice)bool{
	q.Mu.Lock()
	defer q.Mu.Unlock()
	if _,ok:=q.M[notice.Id];ok{
		return false
	}
	q.M[notice.Id]=notice
	q.Enqueue(notice.Id)
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
		pop_id,_:=q.Front()

		if _,ok:=q.M[pop_id];ok{
			if q.M[pop_id].Recvmsg!=nil{
				min_msg:=q.M[pop_id]
				delete(q.M,pop_id)
				q.Dequeue()
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
		pop_id,_:=q.Front()
		if _,ok:=q.M[pop_id];ok{
			if q.M[pop_id].Recvmsg!=nil{
				min_msg:=q.M[pop_id]
				delete(q.M,pop_id)
				q.Dequeue()
				q.Pop<-min_msg
			}
		}
	}
}


func (q *QueueMsg) Length() interface{} {
	return (q.rear - q.front + q.maxSize) % q.maxSize
}

func (q *QueueMsg) Enqueue(e uint16) error {
	if (q.rear+1)%q.maxSize == q.front {
		return fmt.Errorf("quque is full")
	}
	q.Queue[q.rear] = e
	q.rear = (q.rear + 1) % q.maxSize
	return nil
}

func (q *QueueMsg) Dequeue() (e uint16, err error) {
	if q.rear == q.front {
		return e, fmt.Errorf("quque is empty")
	}
	e = q.Queue[q.front]
	q.front = (q.front + 1) % q.maxSize
	return e, nil
}

func (q *QueueMsg) Front() (e uint16, err error) {
	if q.rear == q.front {
		return e, fmt.Errorf("quque is empty")
	}
	e = q.Queue[q.front]
	return e, nil
}