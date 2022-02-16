package bitarrary

import (
	"testing"
)

func TestAnd(t *testing.T) {
	a, b := 0b101, 0b110
	x := FromNumber(a)
	y := FromNumber(b)
	if x.And(y).ToNumber() != uint64(a&b) {
		t.Error("fail")
	}
}
func TestOr(t *testing.T) {
	a, b := 0b101, 0b110
	x := FromNumber(a)
	y := FromNumber(b)
	if x.Or(y).ToNumber() != uint64(a|b) {
		t.Error("fail")
	}
}
func TestXor(t *testing.T) {
	a, b := 0b101, 0b110
	x := FromNumber(a)
	y := FromNumber(b)
	if x.Xor(y).ToNumber() != uint64(a^b) {
		t.Error("fail")
	}
}

func TestNot(t *testing.T) {
	n := uint64(0b101)
	x := FromNumber(n)
	if x.Not().ToNumber() != ^n {
		t.Error("fail")
	}
}

func TestLeftShift(t *testing.T) {
	n := 1
	a := uint64(0b101)
	b := BitSetFromNumber(a)
	b.LShift(n)
	if b.ToNumber()[0] != a<<uint64(n) {
		t.Error("fail")
	}
}

func TestRightShift(t *testing.T) {
	n := 1
	a := uint64(0b101)
	b := FromNumber(a)
	b.RShift(n)
	if b.ToNumber() != a>>uint64(n) {
		t.Error("fail")
	}
}

func TestSub(t *testing.T) {
	a := 5
	b := 1
	if FromNumber(a).Inc(-b).ToNumber() != uint64(a-b) {
		t.Error("fail")
	}
}

func TestAdd(t *testing.T) {
	a := 5
	b := 1
	if FromNumber(a).Inc(b).ToNumber() != uint64(a+b) {
		t.Error("fail")
	}
}

func TestSetBit(t *testing.T) {
	var a = uint16(0b10101)
	x := FromNumber(a)
	x.SetBit(15, true)
	if uint16(x.ToNumber()) != uint16(0b1000000000010101) {
		t.Error("fail")
	}
}

func TestGetBit(t *testing.T) {
	var a = uint16(0b10101010)
	x := FromNumber(a)
	if !x.GetBit(1) {
		t.Error("fail")
	}
}

func TestCut(t *testing.T) {
	var b = NewBitArrary(9)
	b.Reset(true)
	if b.ToNumber() != uint64(0b111111111) {
		t.Error("fail")
	}
}
