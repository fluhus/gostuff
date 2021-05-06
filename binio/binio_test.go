package binio

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestWriteByte(t *testing.T) {
	in := []byte{234, 231, 54, 22, 6, 86}
	buf := bytes.NewBuffer(nil)
	for _, b := range in {
		if err := WriteByte(buf, b); err != nil {
			t.Fatalf("WriteByte(%v)=%q, want success", b, err.Error())
		}
	}
	for _, b := range in {
		got, err := ReadByte(buf)
		if err != nil {
			t.Fatalf("ReadByte()=%q, want success", err.Error())
		}
		if got != b {
			t.Fatalf("ReadByte()=%v, want %v", got, b)
		}
	}
	if _, err := ReadByte(buf); err != io.EOF {
		t.Fatalf("ReadByte()=%q, want EOF", err.Error())
	}
}

func TestWriteUint64(t *testing.T) {
	in := []uint64{23544, 231, 5454, 2762, 665756, 86756}
	buf := bytes.NewBuffer(nil)
	for _, i := range in {
		if err := WriteUint64(buf, i); err != nil {
			t.Fatalf("WriteUint64(%v)=%q, want success", i, err.Error())
		}
	}
	for _, i := range in {
		got, err := ReadUint64(buf)
		if err != nil {
			t.Fatalf("ReadUint64()=%q, want success", err.Error())
		}
		if got != i {
			t.Fatalf("ReadUint64()=%v, want %v", got, i)
		}
	}
	if _, err := ReadUint64(buf); err != io.EOF {
		t.Fatalf("ReadUint64()=%q, want EOF", err.Error())
	}
}

func TestWriteUvarint(t *testing.T) {
	in := []uint64{23544, 231, 5454, 2762, 665756, 86756}
	buf := bytes.NewBuffer(nil)
	for _, i := range in {
		if err := WriteUvarint(buf, i); err != nil {
			t.Fatalf("WriteUvarint(%v)=%q, want success", i, err.Error())
		}
	}
	for _, i := range in {
		got, err := ReadUvarint(buf)
		if err != nil {
			t.Fatalf("ReadUvarint()=%q, want success", err.Error())
		}
		if got != i {
			t.Fatalf("ReadUvarint()=%v, want %v", got, i)
		}
	}
	if _, err := ReadUvarint(buf); err != io.EOF {
		t.Fatalf("ReadUvarint()=%q, want EOF", err.Error())
	}
}

func TestWriteString(t *testing.T) {
	in := []string{"", "amit", "32234312", "fsdfjsd dkjas \" dfsd43.312#@!"}
	buf := bytes.NewBuffer(nil)
	for _, s := range in {
		if err := WriteString(buf, s); err != nil {
			t.Fatalf("WriteString(%v)=%q, want success", s, err.Error())
		}
	}
	for _, s := range in {
		got, err := ReadString(buf)
		if err != nil {
			t.Fatalf("ReadString()=%q, want success", err.Error())
		}
		if got != s {
			t.Fatalf("ReadString()=%v, want %v", got, s)
		}
	}
	if _, err := ReadString(buf); err != io.EOF {
		t.Fatalf("ReadString()=%q, want EOF", err.Error())
	}
}

func TestReadUint64_bad(t *testing.T) {
	in := [][]byte{
		{1},
		{1, 1},
		{1, 1, 1},
		{1, 1, 1, 1},
		{1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1},
	}
	for _, b := range in {
		buf := bytes.NewBuffer(b)
		u, err := ReadUint64(buf)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("ReadUint64(%v)={%v,%v}, want ErrUnexpectedEOF",
				b, u, err)
		}
	}
}

func TestReadString_eof(t *testing.T) {
	in := [][]byte{
		{1},
		{2, 1},
	}
	for _, b := range in {
		buf := bytes.NewBuffer(b)
		s, err := ReadString(buf)
		if err != io.ErrUnexpectedEOF {
			t.Errorf("ReadString(%v)={%v,%v}, want ErrUnexpectedEOF",
				b, s, err)
		}
	}
}

func TestGetBit(t *testing.T) {
	in := []byte{0b10100101, 0b00111100}
	want := []int{1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0}
	var got []int
	for i := 0; i < len(in)*8; i++ {
		got = append(got, GetBit(in, i))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetBit(%v)=%v, want %v", in, got, want)
	}
}

func TestSetBit(t *testing.T) {
	want := []byte{0b10100101, 0b00111100}
	vals := []int{1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0}
	got := []byte{0, 0}
	for i, v := range vals {
		SetBit(got, i, v)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SetBit(%v)=%v, want %v", vals, got, want)
	}
	got = []byte{255, 255}
	for i, v := range vals {
		SetBit(got, i, v)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SetBit(%v)=%v, want %v", vals, got, want)
	}
}

func TestGetHalfByte(t *testing.T) {
	buf := []byte{0b10100101, 0b00111100}
	want := []byte{0b00000101, 0b00001010, 0b00001100, 0b00000011}
	got := []byte{GetHalfByte(buf, 0), GetHalfByte(buf, 1),
		GetHalfByte(buf, 2), GetHalfByte(buf, 3)}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("getHalfByte(%v)=%v, want %v", buf, got, want)
	}
}

func TestSetHalfByte(t *testing.T) {
	got := []byte{0, 0}
	want := []byte{0b10100101, 0b00111100}
	SetHalfByte(got, 0, 0b00000101)
	SetHalfByte(got, 1, 0b00001010)
	SetHalfByte(got, 2, 0b00001100)
	SetHalfByte(got, 3, 0b00000011)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("setHalfByte(...)=%v, want %v", got, want)
	}
}
