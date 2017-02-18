package litebuf

const (
	cacheAlign   = 64 // cache line size
	preallocSize = 32 // preallocSize >= (cacheAlign / 2)

	pageSize = 4096 // page size
	pageMask = pageSize - 1
)

func pageRound(n int) int {
	return (n + pageSize - 1) &^ pageMask
}

type Buffer struct {
	buf []byte
	p   int
	pre [preallocSize]byte
}

func (b *Buffer) Len() int {
	return b.p
}

func (b *Buffer) Cap() int {
	return len(b.buf)
}

func (b *Buffer) Resize(n int) {
	if cap(b.buf) >= n {
		b.buf = b.buf[:n]
		if b.p > n {
			b.p = n
		}
	} else if n <= preallocSize {
		b.buf = b.pre[:]
	} else {
		if n <= pageSize/2 {
			n = (n + n) &^ (cacheAlign - 1)
		} else {
			n = pageRound(n + cap(b.buf)/2)
		}
		buf := make([]byte, n)
		copy(buf, b.buf[:b.p])
		b.buf = buf
	}
}

func (b *Buffer) Reserve(n int) []byte {
	if b.p+n > len(b.buf) {
		b.Resize(b.p + n)
	}
	p := b.p
	b.p += n
	return b.buf[p:b.p]
}

func (b *Buffer) WriteByte(c byte) {
	if b.p >= len(b.buf) {
		b.Resize(b.p + 8)
	}
	b.buf[b.p] = c
	b.p++
}

func (b *Buffer) WriteString(s string) (int, error) {
	if b.p+len(s) > len(b.buf) {
		b.Resize(b.p + len(s))
	}
	n := copy(b.buf[b.p:b.p+len(s)], s)
	b.p += n
	return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	if b.p+len(p) > len(b.buf) {
		b.Resize(b.p + len(p))
	}
	n := copy(b.buf[b.p:b.p+len(p)], p)
	b.p += n
	return n, nil
}

func (b *Buffer) Reset() {
	b.p = 0
	b.buf = b.buf[:cap(b.buf)]
}

func (b *Buffer) Bytes() []byte {
	return b.buf[:b.p]
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}
