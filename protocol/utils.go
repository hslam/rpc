package protocol

func uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

func bytesToUint16(array []byte) uint16 {
	var data uint16 =0
	for i:=0;i< len(array);i++  {
		data = data+uint16(uint(array[i])<<uint(8*i))
	}
	return data
}

func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}

func bytesToUint32(array []byte) uint32 {
	var data uint32 =0
	for i:=0;i< len(array);i++  {
		data = data+uint32(uint(array[i])<<uint(8*i))
	}
	return data
}

func checkSum(b []byte) uint16 {
	sum := 0
	for n := 1; n < len(b)-1; n += 2 {
		sum += int(b[n])<<8 + int(b[n+1])
	}
	lowbit:=sum & 0xffff
	highbit:=sum >> 16
	checksum := lowbit +highbit
	checksum += (checksum >> 16)
	var ans = uint16(^checksum)
	return ans
}