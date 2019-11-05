package protocol

import (
)


const (
	FrameHeaderLength		= 5
)

func PacketFrame(priority uint8,id uint32,body []byte) []byte {
	var buffer=make([]byte,FrameHeaderLength+len(body))
	copy(buffer[0:], []byte{priority})
	copy(buffer[1:], uint32ToBytes(id))
	copy(buffer[FrameHeaderLength:],body)
	return buffer
}

func UnpackFrame(buffer []byte) (priority uint8,id uint32,body []byte,err error) {
	if len(buffer)<FrameHeaderLength{
		return
	}
	priority=buffer[:1][0]
	id = bytesToUint32(buffer[1:FrameHeaderLength])
	body=buffer[FrameHeaderLength:]
	return
}

type Frame struct {
	priority	uint8
	id			uint32
	body		[]byte
}
