package protocol

import (
	"net"
	"fmt"
)

type UDPMsg struct {
	ID uint16
	Data []byte
	RemoteAddr *net.UDPAddr
}

func ReadUDPConn(reader *net.UDPConn, readChan chan *UDPMsg) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	buffer := make([]byte, 65536)
	for {
		func(){
			defer func() {
				if err := recover(); err != nil {
				}
			}()
			n, remoteAddr, err := reader.ReadFromUDP(buffer[0:])
			if err != nil {
				return
			}
			var data =make([]byte,n)
			copy(data,buffer[:n])
			oprationType,id,msg,err:=UnpackMessage(data)
			if err != nil {
				return
			}
			if oprationType==OprationTypeData{
				readChan<-&UDPMsg{id,msg,remoteAddr}
			}
		}()
	}
}

func WriteUDPConn(writer *net.UDPConn, writeChan chan *UDPMsg) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for udp_msg:= range writeChan{
		func(){
			defer func() {
				if err := recover(); err != nil {
				}
			}()
			if udp_msg.Data!=nil{
				fmt.Println()
				_, err := writer.WriteToUDP(PacketMessage(OprationTypeData,udp_msg.ID,udp_msg.Data),udp_msg.RemoteAddr)
				if err != nil {
					return
				}
			}else {
				_, err := writer.WriteToUDP(PacketMessage(OprationTypeAck,udp_msg.ID,nil),udp_msg.RemoteAddr)
				if err != nil {
					return
				}
			}
		}()

	}
}
