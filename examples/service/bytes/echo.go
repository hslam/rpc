package bytes

//Echo defines a value.
type Echo struct {
	value []byte
}

//ToLower returns the lowercase character
func (e *Echo) ToLower(req *[]byte, res *[]byte) error {
	*res = *req
	for i := 0; i < len(*req); i++ {
		if (*req)[i] > 64 && (*req)[i] < 91 {
			(*res)[i] += 32
		}
	}
	return nil
}

//Set sets the value.
func (e *Echo) Set(req *[]byte) error {
	e.value = *req
	return nil
}

//Get gets the value.
func (e *Echo) Get(res *[]byte) error {
	*res = e.value
	return nil
}

//Clear clears the value.
func (e *Echo) Clear() error {
	e.value = nil
	return nil
}
