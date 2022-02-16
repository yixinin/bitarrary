package bitarrary

import (
	"fmt"
	"testing"
)

func TestBitarrary(t *testing.T) {

	var x = FromNumber(uint64(1))
	x.MoveLeft(19)
	fmt.Println(x.ToNumber(), 1<<19)
}

func TestAnd(t *testing.T) {
	x := FromNumber(0b101)
	y := FromNumber(0b110)
	fmt.Println(x.And(y).ToNumber(), 0b101&0b110)
}
func TestOr(t *testing.T) {
	x := FromNumber(0b101)
	y := FromNumber(0b110)
	fmt.Println(x.Or(y).ToNumber(), 0b101|0b110)
}
func TestXor(t *testing.T) {
	x := FromNumber(0b101)
	y := FromNumber(0b110)
	fmt.Println(x.Xor(y).ToNumber(), 0b101^0b110)
}

func TestNot(t *testing.T) {
	n := uint64(0b101)
	x := FromNumber(n)
	fmt.Println(x.Not().ToNumber(), ^n)
}

func TestMoveLeft(t *testing.T) {
	n := uint64(0b101)
	x := FromNumber(n)
	x.MoveLeft(9)
	fmt.Println(x.ToNumber(), n<<9)
}

func TestMoveRight(t *testing.T) {
	n := uint64(0b10100)
	x := FromNumber(n)
	x.MoveRight(2)
	fmt.Println(x.ToNumber(), n>>2)
}

func TestSub(t *testing.T) {
	a := uint16(0)
	fmt.Println(FromNumber(a))
	b := uint16(1)
	x := FromNumber(a - b)
	fmt.Println(x)

	fmt.Println(dec(FromNumber(a).bytes, FromNumber(b).bytes))
}

func TestAdd(t *testing.T) {
	a := uint16(255)
	fmt.Println(FromNumber(a))
	b := uint16(1)
	x := FromNumber(a + b)
	fmt.Println(x)

	fmt.Println(inc(FromNumber(a).bytes, FromNumber(b).bytes))
}

func TestSetBit(t *testing.T) {
	var a = uint16(0b10101)
	x := FromNumber(a)
	x.SetBit(15, true)
	fmt.Println(x.ToNumber(), uint16(0b1000000000010101))
}

func TestGetBit(t *testing.T) {
	var a = uint16(0b10101010)
	x := FromNumber(a)
	for i := 0; i < 16; i++ {
		if x.GetBit(i) {
			fmt.Print(1)
		} else {
			fmt.Print(0)
		}
	}

}
