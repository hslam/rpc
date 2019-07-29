package bytepool

type BytePool struct {
	c chan []byte
	w int
}

func NewBytePool(total int, width int) (bp *BytePool) {
	return &BytePool{
		c: make(chan []byte, total),
		w: width,
	}
}

func (bp *BytePool) Get() (b []byte) {
	select {
	case b = <-bp.c:
	default:
		b = make([]byte, bp.w)
	}
	return
}

func (bp *BytePool) Put(b []byte) {
	select {
	case bp.c <- b:
	default:
	}
}

func (bp *BytePool) Width() (n int) {
	return bp.w
}

func GetMake(width int) (b []byte) {
	b = make([]byte, width)
	return
}