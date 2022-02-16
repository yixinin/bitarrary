package bitarrary

import (
	"encoding/binary"
)

type BitArrary struct {
	bytes    []byte
	bitSize  int
	byteSize int
}

func NewBitArrary(size int) *BitArrary {
	byteSize := size / 8
	if size%8 > 0 {
		byteSize++
	}
	bits := make([]byte, byteSize)

	return &BitArrary{
		bytes:    bits,
		bitSize:  size,
		byteSize: byteSize,
	}
}

func (b *BitArrary) Len() int {
	return b.bitSize
}

func (b *BitArrary) ToNumber() uint64 {
	switch b.byteSize {
	case 1:
		return uint64(b.bytes[0])
	case 2:
		return uint64(binary.BigEndian.Uint16(b.bytes))
	case 4:
		return uint64(binary.BigEndian.Uint32(b.bytes))
	case 8:
		return uint64(binary.BigEndian.Uint64(b.bytes))
	}
	return 0
}

func (b *BitArrary) ToArrary() []bool {
	var bits = make([]bool, 0, b.byteSize*8)
	for i := 0; i < b.byteSize; i++ {
		for j := 7; j >= 0; j-- {
			bits = append(bits, b.bytes[i]&(uint8(1)<<j) > 0)
		}
	}
	return bits[b.byteSize*8-b.bitSize:]
}

func FromNumber(i interface{}) *BitArrary {
	switch v := i.(type) {
	case uint8:
		bits := NewBitArrary(8)
		bits.bytes[0] = v
		return bits
	case uint16:
		bits := NewBitArrary(16)
		binary.BigEndian.PutUint16(bits.bytes, v)
		return bits
	case uint32:
		bits := NewBitArrary(32)
		binary.BigEndian.PutUint32(bits.bytes, v)
		return bits
	case uint64:
		bits := NewBitArrary(64)
		binary.BigEndian.PutUint64(bits.bytes, v)
		return bits
	case uint:
		bits := NewBitArrary(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	case int8:
		bits := NewBitArrary(8)
		bits.bytes[0] = uint8(v)
		return bits
	case int16:
		bits := NewBitArrary(16)
		binary.BigEndian.PutUint16(bits.bytes, uint16(v))
		return bits
	case int32:
		bits := NewBitArrary(32)
		binary.BigEndian.PutUint32(bits.bytes, uint32(v))
		return bits
	case int64:
		bits := NewBitArrary(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	case int:
		bits := NewBitArrary(64)
		binary.BigEndian.PutUint64(bits.bytes, uint64(v))
		return bits
	}
	panic("not a valid number")
}

func (b *BitArrary) Reset(positive bool) *BitArrary {
	for i := range b.bytes {
		if positive {
			b.bytes[i] = 255
		} else {
			b.bytes[i] = 0
		}
	}
	return b.cut()
}

func (b *BitArrary) cut() *BitArrary {
	if b.bitSize < 8*b.byteSize {
		bitSize := b.bitSize % 8
		b.bytes[0] = b.bytes[0] & (uint8(255) >> uint8(8-bitSize))
	}
	return b
}

func (b *BitArrary) Inc(i int) *BitArrary {
	switch i {
	case 0:
		return b
	default:
		if i < 0 {
			x := sub(b.bytes, FromNumber(-i).bytes)
			copy(b.bytes, x[len(x)-b.byteSize:])
		} else {
			x := add(b.bytes, FromNumber(i).bytes)
			copy(b.bytes, x[len(x)-b.byteSize:])
		}
	}
	return b.cut()
}

func add(a, b []byte) []byte {
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
		} else {
			over = 0
		}
		ret[i] = o
		ai--
		bi--
	}
	return ret
}

func sub(a, b []byte) []byte {
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
		} else {
			take = 0
		}
		ret[i] = o
		ai--
		bi--
	}
	return ret
}

func (b *BitArrary) SetBit(index int, positive bool) *BitArrary {
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
	return b.cut()
}

func (b *BitArrary) GetBit(index int) bool {
	idx := index / 8
	if idx >= b.byteSize {
		return false
	}
	idx = b.byteSize - idx - 1
	cidx := index % 8
	return b.bytes[idx]&(uint8(0b1)<<byte(cidx)) > 0
}

func (b *BitArrary) And(other *BitArrary) *BitArrary {
	return And(b, other)
}

func (b *BitArrary) Or(other *BitArrary) *BitArrary {
	return Or(b, other)
}

func (b *BitArrary) Xor(other *BitArrary) *BitArrary {
	return Xor(b, other)
}

func (b *BitArrary) Not() *BitArrary {
	return Not(b)
}

func (b *BitArrary) LShift(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.RShift(-n)
		return
	}
	lshift(b.bytes, n)
	b.cut()
}

func lshift(b []byte, n int) {
	size := len(b)
	bitSize := size * 8

	var dstByteIdx, dstBitOff int
	var srcByteIdx, dstBitIdx, srcBitOff int

	for i := n; i < bitSize; i++ {
		srcByteIdx = i / 8
		srcBitOff = i % 8

		dstBitIdx = i - n
		dstByteIdx = dstBitIdx / 8
		dstBitOff = dstBitIdx % 8

		srcBitVal := b[srcByteIdx] & (uint8(1) << uint8(7-srcBitOff))
		if srcBitVal > 0 {
			b[dstByteIdx] = b[dstByteIdx] | (uint8(1) << uint8(7-dstBitOff))
		} else {
			b[dstByteIdx] = b[dstByteIdx] & (^(uint8(1) << uint8(7-dstBitOff)))
		}
	}

	b[dstByteIdx] = b[dstByteIdx] & (uint8(255) << (7 - dstBitOff))
	for i := dstByteIdx + 1; i < size; i++ {
		b[i] = 0
	}
}

func (b *BitArrary) RShift(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.LShift(-n)
		return
	}
	rshift(b.bytes, n)
	b.cut()
}

func rshift(b []byte, n int) {
	size := len(b)
	bitSize := size * 8

	var dstByteIdx, dstBitOff int
	var srcByteIdx, dstBitIdx, srcBitOff int
	for i := bitSize - n - 1; i >= 0; i-- {
		srcByteIdx = i / 8
		srcBitOff = i % 8

		dstBitIdx = i + n
		dstByteIdx = dstBitIdx / 8
		dstBitOff = dstBitIdx % 8

		srcBitVal := b[srcByteIdx] & (uint8(1) << uint8(7-srcBitOff))
		if srcBitVal > 0 {
			b[dstByteIdx] = b[dstByteIdx] | (uint8(1) << uint8(7-dstBitOff))
		} else {
			b[dstByteIdx] = b[dstByteIdx] & (^(uint8(1) << uint8(7-dstBitOff)))
		}
	}

	b[dstByteIdx] = b[dstByteIdx] & (uint8(255) << (7 - dstBitOff))
	for i := dstByteIdx - 1; i >= 0; i-- {
		b[i] = 0
	}
}

func (b *BitArrary) Compare(other *BitArrary) int {
	return Compare(b, other)
}

func And(a, b *BitArrary) *BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb < size {
		size = lb
	}
	var ret = NewBitArrary(size * 8)
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		ret.bytes[i] = a.bytes[ai] & b.bytes[bi]
		ai--
		bi--
	}
	return ret.cut()
}

func Or(a, b *BitArrary) *BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitArrary(size * 8)
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
	return ret.cut()
}

func Xor(a, b *BitArrary) *BitArrary {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitArrary(size * 8)
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
	return ret.cut()
}

func Not(b *BitArrary) *BitArrary {
	var ret = NewBitArrary(b.bitSize)
	for i := range b.bytes {
		ret.bytes[i] = ^b.bytes[i]
	}
	return ret.cut()
}

// return -1 if a<b
// return 1 if a>b
// return 0 if a==b
func Compare(a, b *BitArrary) int {
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
			return -1
		} else if vb < va {
			return 1
		}
		ai--
		bi--
	}
	return 0
}
