package flu

import "io"

type Counter int64

func (c *Counter) Value() int64 {
	return *(*int64)(c)
}

func (c *Counter) Add(n int64) {
	*(*int64)(c) += n
}

type ReaderCounter struct {
	io.Reader
	*Counter
}

func (rc ReaderCounter) Read(data []byte) (int, error) {
	n, err := rc.Reader.Read(data)
	if err != nil {
		return 0, err
	}

	rc.Add(int64(n))
	return n, nil
}

func (rc ReaderCounter) Close() error {
	return ReaderCloser{rc.Reader}.Close()
}

type WriterCounter struct {
	io.Writer
	*Counter
}

func (wc WriterCounter) Write(data []byte) (int, error) {
	n, err := wc.Writer.Write(data)
	if err != nil {
		return 0, err
	}

	wc.Add(int64(n))
	return n, nil
}

func (wc WriterCounter) Close() error {
	return WriterCloser{wc.Writer}.Close()
}

type IOCounter struct {
	Input
	Output
	Counter
}

func (c *IOCounter) Reader() (io.Reader, error) {
	r, err := c.Input.Reader()
	if err != nil {
		return nil, err
	}
	return ReaderCounter{Reader: r, Counter: &c.Counter}, nil
}

func (c *IOCounter) Writer() (io.Writer, error) {
	w, err := c.Output.Writer()
	if err != nil {
		return nil, err
	}
	return WriterCounter{Writer: w, Counter: &c.Counter}, nil
}
