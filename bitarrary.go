package bitarrary

import (
	"encoding/binary"
)

type BitArrary struct {
	bytes    []byte
	bitSize  int
	byteSize int
}

func NewBitArray(size int) BitArrary {
	byteSize := size / 8
	if size%8 > 0 {
		byteSize++
	}
	bits := make([]byte, byteSize)

	return BitArrary{
		bytes:    bits,
		bitSize:  size,
		byteSize: byteSize,
	}
}

func (b BitArrary) ToNumber() interface{} {
	switch b.byteSize {
	case 1:
		return b.bytes[0]
	case 2:
		return binary.BigEndian.Uint16(b.bytes)
	case 4:
		return binary.BigEndian.Uint32(b.bytes)
	case 8:
		return binary.BigEndian.Uint64(b.bytes)
	}
	return -1
}

func FromNumber(i interface{}) BitArrary {
	switch v := i.(type) {
	case uint8:
		bits := NewBitArray(8)
		bits.bytes[0] = v
		return bits
	case uint16:
		bits := NewBitArray(16)
		binary.BigEndian.PutUint16(bits.bytes, v)
		return bits
	case uint32:
		bits := NewBitArray(32)
		binary.BigEndian.PutUint32(bits.bytes, v)
		return bits
	case uint64:
		bits := NewBitArray(64)
		binary.BigEndian.PutUint64(bits.bytes, v)
		return bits
	case uint:
		bits := NewBitArray(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	case int8:
		bits := NewBitArray(8)
		bits.bytes[0] = uint8(v)
		return bits
	case int16:
		bits := NewBitArray(16)
		binary.BigEndian.PutUint16(bits.bytes, uint16(v))
		return bits
	case int32:
		bits := NewBitArray(32)
		binary.BigEndian.PutUint32(bits.bytes, uint32(v))
		return bits
	case int64:
		bits := NewBitArray(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	case int:
		bits := NewBitArray(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	}
	panic("not surport this number")
}

func (b BitArrary) Reset(positive bool) {
	for i := range b.bytes {
		if positive {
			b.bytes[i] = 1
		} else {
			b.bytes[i] = 0
		}
	}
}

func (b BitArrary) Inc(i int) {
	switch i {
	case 0:
		return
	default:
		if i < 0 {
			x := dec(b.bytes, FromNumber(-i).bytes)
			copy(b.bytes, x[:len(b.bytes)])
		} else {
			x := inc(b.bytes, FromNumber(i).bytes)
			copy(b.bytes, x[:len(b.bytes)])
		}
	}
}

func inc(a, b []byte) []byte {
	var size = len(a)
	lb := len(b)
	if lb > size {
		size = lb
	}
	ai, bi := len(a)-1, len(b)-1

	ret := make([]byte, size)

	over := byte(0)
	for i := size - 1; i >= 0; i-- {
		va, vb := byte(0), byte(0)
		if ai >= 0 {
			va = a[ai]
		}
		if bi >= 0 {
			vb = b[bi]
		}

		o := va + vb + over
		if o < va {
			over = 1
			ret[i] = 0
		} else {
			ret[i] = o
		}
		ai--
		bi--
	}
	return ret
}

func dec(a, b []byte) []byte {
	var size = len(a)
	lb := len(b)
	if lb > size {
		size = lb
	}
	ai, bi := len(a)-1, len(b)-1

	ret := make([]byte, size)

	take := byte(0)
	for i := size - 1; i >= 0; i-- {
		va, vb := byte(0), byte(0)
		if ai >= 0 {
			va = a[ai]
		}
		if bi >= 0 {
			vb = b[bi]
		}

		o := va - take - vb
		if o > va {
			take = 1
			ret[i] = o
		} else {
			ret[i] = o
		}
		ai--
		bi--
	}
	return ret
}

func (b BitArrary) SetBit(index int, positive bool) {
	bidx := index / 8
	if bidx >= b.byteSize {
		panic("over flow")
	}
	bidx = b.byteSize - bidx - 1

	bbidx := index % 8
	if positive {
		b.bytes[bidx] = b.bytes[bidx] | uint8(0b1)<<byte(bbidx)
	} else {
		b.bytes[bidx] = b.bytes[bidx] & ^(uint8(0b1) << byte(bbidx))
	}
}

func (b BitArrary) GetBit(index int) bool {
	idx := index / 8
	if idx >= b.byteSize {
		return false
	}
	idx = b.byteSize - idx - 1
	cidx := index % 8
	return b.bytes[idx]&(uint8(0b1)<<byte(cidx)) > 0
}

func (b BitArrary) And(other BitArrary) BitArrary {
	return And(b, other)
}

func (b BitArrary) Or(other BitArrary) BitArrary {
	return Or(b, other)
}

func (b BitArrary) Xor(other BitArrary) BitArrary {
	return Xor(b, other)
}

func (b BitArrary) Not() BitArrary {
	return Not(b)
}

func (b BitArrary) MoveLeft(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.MoveRight(-n)
		return
	}
	for i := 0; i < n; i++ {
		b.moveLeft()
	}
}

func moveLeft(b []byte, i int) {
	// var size = len(b)
	// var bitSize = size * 8

	// minIdx := i / 8
	// maxIdx := size - minIdx
	// for i := size - 1; i >= 0; i-- {

	// }
}

func (b BitArrary) moveLeft() {
	for i := 0; i < b.byteSize-1; i++ {
		b.bytes[i] = b.bytes[i] << 1
		if b.bytes[i+1]&uint8(0b10000000) > 0 {
			b.bytes[i] = b.bytes[i] | uint8(0b1)
		}
	}
	b.bytes[b.byteSize-1] = b.bytes[b.byteSize-1] << 1
}

func (b BitArrary) MoveRight(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.MoveLeft(-n)
		return
	}
	for i := 0; i < n; i++ {
		b.moveRight()
	}
}

func (b BitArrary) moveRight() {
	for i := b.byteSize - 1; i > 0; i-- {
		b.bytes[i] = b.bytes[i] >> 1
		if b.bytes[i-1]&uint8(0b1) > 0 {
			b.bytes[i] = b.bytes[i] | uint8(0b10000000)
		}
	}
	b.bytes[0] = b.bytes[0] >> 1
}

func (b BitArrary) Compare(other BitArrary) int {
	return Compare(b, other)
}

func And(a, b BitArrary) BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb < size {
		size = lb
	}
	var ret = NewBitArray(size * 8)
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		ret.bytes[i] = a.bytes[ai] & b.bytes[bi]
		ai--
		bi--
	}
	return ret
}

func Or(a, b BitArrary) BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitArray(size * 8)
	for i := size - 1; i >= 0; i-- {
		va, vb := byte(0), byte(0)
		if ai >= 0 {
			va = a.bytes[ai]
		}
		if bi >= 0 {
			vb = b.bytes[bi]
		}
		ret.bytes[i] = va | vb
		ai--
		bi--
	}
	return ret
}

func Xor(a, b BitArrary) BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitArray(size * 8)
	for i := size - 1; i >= 0; i-- {
		va, vb := byte(0), byte(0)
		if ai >= 0 {
			va = a.bytes[ai]
		}
		if bi >= 0 {
			vb = b.bytes[bi]
		}
		ret.bytes[i] = va ^ vb
		ai--
		bi--
	}
	return ret
}

func Not(b BitArrary) BitArrary {
	var ret = NewBitArray(b.bitSize)
	for i := range b.bytes {
		ret.bytes[i] = ^b.bytes[i]
	}
	return ret
}

func Compare(a, b BitArrary) int {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		va, vb := byte(0), byte(0)
		if ai >= 0 {
			va = a.bytes[ai]
		}
		if bi >= 0 {
			vb = b.bytes[bi]
		}
		if vb > va {
			return 1
		} else if vb < va {
			return -1
		}
		ai--
		bi--
	}
	return 0
}
