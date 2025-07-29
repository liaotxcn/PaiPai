package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size == 0 {
		size = 250
	}

	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

func (b *Bitmap) Set(id string) {
	// 计算id在哪个bit
	idx := hash(id) % b.size
	// 计算在哪个byte
	byteIdx := idx / 8
	// 计算在byte的哪个bit位置
	bitIdx := idx % 8

	b.bits[byteIdx] |= 1 << (bitIdx)
}

func hash(id string) int {
	// BKDR哈希算法
	seed := 131313
	hash := 0
	for _, c := range id {
		hash = hash * seed * int(c)
	}
	return hash & 0x7FFFFFF
}

func (b Bitmap) IsSet(id string) bool {
	// 计算id在哪个bit
	idx := hash(id) % b.size
	// 计算在哪个byte
	byteIdx := idx / 8
	// 计算在byte的哪个bit位置
	bitIdx := idx % 8
	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

func Lood(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}
