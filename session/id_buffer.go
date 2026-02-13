package session

import "io"

type idBuffer struct {
	buf []byte
	n   int
}

func newIdBuffer(length int) *idBuffer {
	return &idBuffer{
		buf: make([]byte, length),
	}
}

func (i *idBuffer) ReadString(s string) int {
	n := copy(i.buf[i.n:], s)
	i.n += n

	return n
}

func (i *idBuffer) Read(reader io.Reader) (int, error) {
	n, err := reader.Read(i.buf[i.n:])
	i.n += n

	return n, err
}

func (i *idBuffer) Bytes() []byte {
	return i.buf[:i.n]
}

func (i *idBuffer) Rewind() {
	i.n = 0
}
