package bitarrary

import (
	"encoding/binary"
)

type BitArrary struct {
	bytes    []byte
	bitSize  int
	byteSize int
}

func NewBitArray(size int) *BitArrary {
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

func (b *BitArrary) ToArrary() []uint8 {
	var bits = make([]byte, 0, b.bitSize)
	for i := 0; i < b.byteSize; i++ {
		for j := 7; j >= 0; j-- {
			if b.bytes[i]&(uint8(1)<<j) > 0 {
				bits = append(bits, 1)
			} else {
				bits = append(bits, 0)
			}
		}
	}
	return bits
}

func FromNumber(i interface{}) *BitArrary {
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

func (b *BitArrary) Reset(positive bool) {
	for i := range b.bytes {
		if positive {
			b.bytes[i] = 255
		} else {
			b.bytes[i] = 0
		}
	}
	b.Cut()
}

func (b *BitArrary) Cut() *BitArrary {
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
			copy(b.bytes, x[:len(b.bytes)])
		} else {
			x := add(b.bytes, FromNumber(i).bytes)
			copy(b.bytes, x[:len(b.bytes)])
		}
	}
	return b.Cut()
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
			ret[i] = 0
		} else {
			ret[i] = o
		}
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
			ret[i] = o
		} else {
			ret[i] = o
		}
		ai--
		bi--
	}
	return ret
}

func (b *BitArrary) SetBit(index int, positive bool) {
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

func (b *BitArrary) MoveLeft(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.MoveRight(-n)
		return
	}
	moveLeft(b.bytes, n)
	b.Cut()
}

func moveLeft(b []byte, offset int) {
	size := len(b)
	bitSize := size * 8

	var dstByteIdx, dstBitOff int
	var srcByteIdx, dstBitIdx, srcBitOff int

	for i := offset; i < bitSize; i++ {
		srcByteIdx = i / 8
		srcBitOff = i % 8

		dstBitIdx = i - offset
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

func (b *BitArrary) MoveRight(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.MoveLeft(-n)
		return
	}
	moveRight(b.bytes, n)
	b.Cut()
}

func moveRight(b []byte, n int) {
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
	var ret = NewBitArray(size * 8)
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		ret.bytes[i] = a.bytes[ai] & b.bytes[bi]
		ai--
		bi--
	}
	return ret.Cut()
}

func Or(a, b *BitArrary) *BitArrary {
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
	return ret.Cut()
}

func Xor(a, b *BitArrary) *BitArrary {
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
	return ret.Cut()
}

func Not(b *BitArrary) *BitArrary {
	var ret = NewBitArray(b.bitSize)
	for i := range b.bytes {
		ret.bytes[i] = ^b.bytes[i]
	}
	return ret.Cut()
}

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
			return 1
		} else if vb < va {
			return -1
		}
		ai--
		bi--
	}
	return 0
}
