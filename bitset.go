package bitarrary

type BitSet struct {
	sets     []uint64
	bitSize  int
	byteSize int
}

func NewBitSet(size int) *BitSet {
	byteSize := size / 64
	if size%64 > 0 {
		byteSize++
	}
	bits := make([]uint64, byteSize)

	return &BitSet{
		sets:     bits,
		bitSize:  size,
		byteSize: byteSize,
	}
}

func (b *BitSet) Len() int {
	return b.bitSize
}

func (b *BitSet) ToNumber() []uint64 {
	return b.sets
}

func BitSetFromNumber(i interface{}) *BitSet {
	switch v := i.(type) {
	case uint8:
		bits := NewBitSet(8)
		bits.sets[0] = uint64(v)
		return bits
	case uint16:
		bits := NewBitSet(16)
		bits.sets[0] = uint64(v)
		return bits
	case uint32:
		bits := NewBitSet(32)
		bits.sets[0] = uint64(v)
		return bits
	case uint64:
		bits := NewBitSet(64)
		bits.sets[0] = uint64(v)
		return bits
	case uint:
		bits := NewBitSet(64)
		bits.sets[0] = uint64(v)
		return bits
	case int8:
		bits := NewBitSet(8)
		bits.sets[0] = uint64(v)
		return bits
	case int16:
		bits := NewBitSet(16)
		bits.sets[0] = uint64(v)
		return bits
	case int32:
		bits := NewBitSet(32)
		bits.sets[0] = uint64(v)
		return bits
	case int64:
		bits := NewBitSet(64)
		bits.sets[0] = uint64(v)
		return bits
	case int:
		bits := NewBitSet(64)
		bits.sets[0] = uint64(v)
		return bits
	}
	panic("not a valid number")
}

func (b *BitSet) ToArrary() []bool {
	var bits = make([]bool, 0, b.byteSize*64)
	for i := 0; i < b.byteSize; i++ {
		for j := 63; j >= 0; j-- {
			bits = append(bits, b.sets[i]&(uint64(1)<<j) > 0)
		}
	}
	return bits[b.byteSize*64-b.bitSize:]
}

func (b *BitSet) Reset(positive bool) *BitSet {
	for i := range b.sets {
		if positive {
			b.sets[i] = ^(uint64(0))
		} else {
			b.sets[i] = 0
		}
	}
	return b.cut()
}

func (b *BitSet) cut() *BitSet {
	if b.bitSize < 8*b.byteSize {
		bitSize := b.bitSize % 8
		b.sets[0] = b.sets[0] & ((^uint64(0)) >> uint64(64-bitSize))
	}
	return b
}

func (b *BitSet) Inc(i int) *BitSet {
	switch i {
	case 0:
		return b
	default:
		if i < 0 {
			x := setSub(b.sets, []uint64{uint64(-i)})
			copy(b.sets, x[len(x)-b.byteSize:])
		} else {
			x := setAdd(b.sets, []uint64{uint64(i)})
			copy(b.sets, x[len(x)-b.byteSize:])
		}
	}
	return b.cut()
}

func setAdd(a, b []uint64) []uint64 {
	var size = len(a)
	lb := len(b)
	if lb > size {
		size = lb
	}
	ai, bi := len(a)-1, len(b)-1

	ret := make([]uint64, size)

	over := uint64(0)
	for i := size - 1; i >= 0; i-- {
		va, vb := uint64(0), uint64(0)
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

func setSub(a, b []uint64) []uint64 {
	var size = len(a)
	lb := len(b)
	if lb > size {
		size = lb
	}
	ai, bi := len(a)-1, len(b)-1

	ret := make([]uint64, size)

	take := uint64(0)
	for i := size - 1; i >= 0; i-- {
		va, vb := uint64(0), uint64(0)
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

func (b *BitSet) SetBit(index int, positive bool) *BitSet {
	bidx := index / 64
	if bidx >= b.byteSize {
		panic("over flow")
	}
	bidx = b.byteSize - bidx - 1

	bbidx := index % 64
	if positive {
		b.sets[bidx] = b.sets[bidx] | uint64(0b1)<<uint64(bbidx)
	} else {
		b.sets[bidx] = b.sets[bidx] & ^(uint64(0b1) << uint64(bbidx))
	}
	return b.cut()
}

func (b *BitSet) GetBit(index int) bool {
	idx := index / 64
	if idx >= b.byteSize {
		return false
	}
	idx = b.byteSize - idx - 1
	cidx := index % 64
	return b.sets[idx]&(uint64(0b1)<<uint64(cidx)) > 0
}

func (b *BitSet) And(other *BitSet) *BitSet {
	return setAnd(b, other)
}

func (b *BitSet) Or(other *BitSet) *BitSet {
	return setOr(b, other)
}

func (b *BitSet) Xor(other *BitSet) *BitSet {
	return setXor(b, other)
}

func (b *BitSet) Not() *BitSet {
	return setNot(b)
}

func (b *BitSet) LShift(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.RShift(-n)
		return
	}
	setLshift(b.sets, n)
	b.cut()
}

func setLshift(b []uint64, n int) {
	size := len(b)
	bitSize := size * 64

	var dstByteIdx, dstBitOff int
	var srcByteIdx, dstBitIdx, srcBitOff int

	for i := n; i < bitSize; i++ {
		srcByteIdx = i / 64
		srcBitOff = i % 64

		dstBitIdx = i - n
		dstByteIdx = dstBitIdx / 64
		dstBitOff = dstBitIdx % 64

		srcBitVal := b[srcByteIdx] & (uint64(1) << uint64(63-srcBitOff))
		if srcBitVal > 0 {
			b[dstByteIdx] = b[dstByteIdx] | (uint64(1) << uint64(63-dstBitOff))
		} else {
			b[dstByteIdx] = b[dstByteIdx] & (^(uint64(1) << uint64(63-dstBitOff)))
		}
	}

	b[dstByteIdx] = b[dstByteIdx] & (uint64(255) << (63 - dstBitOff))
	for i := dstByteIdx + 1; i < size; i++ {
		b[i] = 0
	}
}

func (b *BitSet) RShift(n int) {
	if n == 0 {
		return
	}
	if n < 0 {
		b.LShift(-n)
		return
	}
	setRshift(b.sets, n)
	b.cut()
}

func setRshift(b []uint64, n int) {
	size := len(b)
	bitSize := size * 8

	var dstByteIdx, dstBitOff int
	var srcByteIdx, dstBitIdx, srcBitOff int
	for i := bitSize - n - 1; i >= 0; i-- {
		srcByteIdx = i / 64
		srcBitOff = i % 64

		dstBitIdx = i + n
		dstByteIdx = dstBitIdx / 64
		dstBitOff = dstBitIdx % 64

		srcBitVal := b[srcByteIdx] & (uint64(1) << uint64(63-srcBitOff))
		if srcBitVal > 0 {
			b[dstByteIdx] = b[dstByteIdx] | (uint64(1) << uint64(63-dstBitOff))
		} else {
			b[dstByteIdx] = b[dstByteIdx] & (^(uint64(1) << uint64(63-dstBitOff)))
		}
	}

	b[dstByteIdx] = b[dstByteIdx] & ((^uint64(0)) << (63 - dstBitOff))
	for i := dstByteIdx - 1; i >= 0; i-- {
		b[i] = 0
	}
}

func (b *BitSet) Compare(other *BitSet) int {
	return setCompare(b, other)
}

func setAnd(a, b *BitSet) *BitSet {
	var size = a.byteSize
	lb := b.byteSize
	if lb < size {
		size = lb
	}
	var ret = NewBitSet(size * 64)
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		ret.sets[i] = a.sets[ai] & b.sets[bi]
		ai--
		bi--
	}
	return ret.cut()
}

func setOr(a, b *BitSet) *BitSet {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitSet(size * 64)
	for i := size - 1; i >= 0; i-- {
		va, vb := uint64(0), uint64(0)
		if ai >= 0 {
			va = a.sets[ai]
		}
		if bi >= 0 {
			vb = b.sets[bi]
		}
		ret.sets[i] = va | vb
		ai--
		bi--
	}
	return ret.cut()
}

func setXor(a, b *BitSet) *BitSet {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	var ret = NewBitSet(size * 64)
	for i := size - 1; i >= 0; i-- {
		va, vb := uint64(0), uint64(0)
		if ai >= 0 {
			va = a.sets[ai]
		}
		if bi >= 0 {
			vb = b.sets[bi]
		}
		ret.sets[i] = va ^ vb
		ai--
		bi--
	}
	return ret.cut()
}

func setNot(b *BitSet) *BitSet {
	var ret = NewBitSet(b.bitSize)
	for i := range b.sets {
		ret.sets[i] = ^b.sets[i]
	}
	return ret.cut()
}

// return -1 if a<b
// return 1 if a>b
// return 0 if a==b
func setCompare(a, b *BitSet) int {
	var size = a.byteSize
	lb := b.byteSize
	if lb > size {
		size = lb
	}
	var ai = a.byteSize - 1
	var bi = b.byteSize - 1
	for i := size - 1; i >= 0; i-- {
		va, vb := uint64(0), uint64(0)
		if ai >= 0 {
			va = a.sets[ai]
		}
		if bi >= 0 {
			vb = b.sets[bi]
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
