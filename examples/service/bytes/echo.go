package bytes

type Echo struct {
	value []byte
}

func (this *Echo) ToLower(req *[]byte, res *[]byte) error {
	*res = *req
	for i := 0; i < len(*req); i++ {
		if (*req)[i] > 64 && (*req)[i] < 91 {
			(*res)[i] += 32
		}
	}
	return nil
}

func (this *Echo) Set(req *[]byte) error {
	this.value = *req
	return nil
}

func (this *Echo) Get(res *[]byte) error {
	*res = this.value
	return nil
}
func (this *Echo) Clear() error {
	this.value = nil
	return nil
}
