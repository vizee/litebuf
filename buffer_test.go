package litebuf

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBufferStruct(t *testing.T) {
	rt := reflect.TypeOf(Buffer{})
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		t.Logf("%-8s %4d %4d", f.Name, f.Offset, f.Type.Size())
	}
}

func TestWrite(t *testing.T) {
	buf := new(Buffer)
	t.Log(buf.Len(), buf.Cap())

	// WriteString
	buf.WriteString("ab")
	t.Log(buf.Len(), buf.Cap())

	// Write
	buf.Write([]byte{'c', 'd', 'e'})
	t.Log(buf.Len(), buf.Cap())

	buf.WriteString("fghij")
	t.Log(buf.Len(), buf.Cap())

	// WriteByte
	buf.WriteByte('k')
	t.Log(buf.Len(), buf.Cap())

	// Reserve
	resv := buf.Reserve(buf.Cap() - buf.Len() + 1)
	t.Log("reserve", len(resv))
	for i := 0; i < len(resv); i++ {
		resv[i] = 'A' + byte(i)
	}
	t.Log(buf.Len(), buf.Cap())

	// Stringer
	t.Log(buf)

	// valid
	should := "abcdefghijkABCDEFGHIJKLMNOPQRSTUV"
	if !bytes.Equal(buf.Bytes(), []byte(should)) {
		t.Fatal("valid failed")
	}

	buf.Resize(0)
	t.Log(buf.Len(), buf.Cap())

	buf.Reset()
	t.Log(buf.Len(), buf.Cap())
}

func BenchmarkGoBuffer(b *testing.B) {
	b.ReportAllocs()
	buf := bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Write([]byte{'a', 'b', 'c', 'd', 'e'})
	}
	b.Log(buf.Len(), buf.Cap())
}

func BenchmarkLiteBuffer(b *testing.B) {
	b.ReportAllocs()
	buf := Buffer{}
	for i := 0; i < b.N; i++ {
		buf.Write([]byte{'a', 'b', 'c', 'd', 'e'})
	}
	b.Log(buf.Len(), buf.Cap())
}
