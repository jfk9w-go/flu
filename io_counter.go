package flu

type IOCounter int64

func (c *IOCounter) Write(data []byte) (n int, err error) {
	n = len(data)
	c.Add(int64(n))
	return n, nil
}

func (c *IOCounter) Count(encoder EncoderTo) error {
	return EncodeTo(encoder, IO{W: c})
}

func (c *IOCounter) Value() int64 {
	return *(*int64)(c)
}

func (c *IOCounter) Add(n int64) {
	*(*int64)(c) += n
}
