package protocol

import (
	"net"
)

type UDPMsg struct {
	ID uint16
	Data []byte
	RemoteAddr *net.UDPAddr
}

func ReadUDPConn(reader *net.UDPConn, readChan chan *UDPMsg, stopChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buffer := make([]byte, 65536)
	for {
		n, remoteAddr, err := reader.ReadFromUDP(buffer[0:])
		if err != nil {
			continue
		}
		var data =make([]byte,n)
		copy(data,buffer[:n])
		oprationType,id,msg,err:=UnpackMessage(data)
		if err != nil {
			continue
		}
		if oprationType==OprationTypeData{
			readChan<-&UDPMsg{id,msg,remoteAddr}
		}
	}
	stopChan <- true
}

func WriteUDPConn(writer *net.UDPConn, writeChan chan *UDPMsg, stopChan chan bool) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for udp_msg:= range writeChan{
		if udp_msg.Data!=nil{
			_, err := writer.WriteToUDP(PacketMessage(OprationTypeData,udp_msg.ID,udp_msg.Data),udp_msg.RemoteAddr)
			if err != nil {
				continue
			}
		}else {
			_, err := writer.WriteToUDP(PacketMessage(OprationTypeAck,udp_msg.ID,nil),udp_msg.RemoteAddr)
			if err != nil {
				continue
			}
		}
	}
}
