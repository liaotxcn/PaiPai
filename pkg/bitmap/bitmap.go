package bitmap

import (
	"errors"
	"sync"
)

// 常量定义
const (
	defaultSize = 250     // 默认位图大小（字节），可存储2000个bit（250*8）
	maxSize     = 1 << 28 // 最大支持256MB大小的位图（2^28 bytes）
)

// 错误定义
var (
	ErrIndexOutOfRange = errors.New("bit position exceeds bitmap capacity")
	ErrInvalidSize     = errors.New("invalid size parameter (must be positive and <= maxSize)")
)

// Bitmap 核心数据结构
type Bitmap struct {
	bits    []byte       // 底层存储，每个byte存储8个bit
	size    int          // 位图总容量（bit数 = len(bits)*8）
	version uint8        // 版本号，用于跟踪修改（Resize时递增）
	mu      sync.RWMutex // 读写锁，保证并发安全
}

/****************************** 初始化方法 ******************************/

// NewBitmap 创建新位图实例
// 参数：
//   - size: 位图字节大小，0表示使用默认大小
//
// 返回值：
//   - *Bitmap: 新位图实例
//   - error: 参数非法时返回错误
func NewBitmap(size int) (*Bitmap, error) {
	// 参数校验
	if size < 0 {
		return nil, ErrInvalidSize
	}
	if size == 0 {
		size = defaultSize // 使用默认大小
	}
	if size > maxSize {
		return nil, ErrInvalidSize
	}

	// 初始化位图
	return &Bitmap{
		bits:    make([]byte, size), // 分配内存
		size:    size * 8,           // 计算总bit数
		version: 1,                  // 初始版本号
	}, nil
}

/****************************** 基础操作 ******************************/

// Set 通过字符串ID设置对应bit为1（线程安全）
// 实现步骤：
//  1. 计算字符串哈希值确定bit位置
//  2. 调用SetBit执行实际设置操作
func (b *Bitmap) Set(id string) error {
	idx := b.hash(id)
	return b.SetBit(idx)
}

// SetBit 直接设置指定bit位置为1（线程安全）
// 关键逻辑：
//   - 使用位操作快速定位到具体bit
//   - 线程安全写锁保护
func (b *Bitmap) SetBit(pos int) error {
	// 边界检查
	if pos < 0 || pos >= b.size {
		return ErrIndexOutOfRange
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// 计算字节位置和bit偏移量
	byteIdx := pos / 8      // 确定哪个byte存储目标bit
	bitIdx := uint(pos % 8) // 确定byte内的具体bit位置（0-7）

	// 使用位操作设置bit（OR操作）
	b.bits[byteIdx] |= 1 << bitIdx
	return nil
}

// IsSet 检查字符串ID对应的bit是否为1（线程安全）
func (b *Bitmap) IsSet(id string) (bool, error) {
	idx := b.hash(id)
	return b.IsBitSet(idx)
}

// IsBitSet 直接检查指定bit位置（线程安全）
// 关键逻辑：
//   - 使用位掩码快速检查bit状态
//   - 读锁保护并发读取
func (b *Bitmap) IsBitSet(pos int) (bool, error) {
	if pos < 0 || pos >= b.size {
		return false, ErrIndexOutOfRange
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	byteIdx := pos / 8
	bitIdx := uint(pos % 8)
	// 使用AND操作检查bit是否为1
	return (b.bits[byteIdx] & (1 << bitIdx)) != 0, nil
}

func (b *Bitmap) Clear(id string) error {
	idx := b.hash(id)
	return b.ClearBit(idx)
}

func (b *Bitmap) ClearBit(pos int) error {
	if pos < 0 || pos >= b.size {
		return ErrIndexOutOfRange
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	byteIdx := pos / 8
	bitIdx := uint(pos % 8)
	b.bits[byteIdx] &^= 1 << bitIdx
	return nil
}

func (b *Bitmap) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := range b.bits {
		b.bits[i] = 0
	}
}

/****************************** 高级操作 ******************************/

// Count 统计位图中设置为1的bit数量（线程安全）
// 优化点：
//   - 使用Brian Kernighan算法高效计算1的个数
//   - 时间复杂度：O(n)其中n为1的个数（优于遍历每个bit）
func (b *Bitmap) Count() int {
	count := 0
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, byteVal := range b.bits {
		// Brian Kernighan算法核心：
		// 每次清除最低位的1，直到byte变为0
		for byteVal != 0 {
			byteVal &= byteVal - 1 // 清除最低位的1
			count++
		}
	}
	return count
}

// Resize 动态调整位图大小（线程安全）
// 关键逻辑：
//  1. 创建新大小的字节数组
//  2. 复制原有数据（自动截断或补零）
//  3. 更新版本号标识修改
func (b *Bitmap) Resize(newSize int) error {
	// 参数校验
	if newSize <= 0 || newSize > maxSize {
		return ErrInvalidSize
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	// 大小未变化时直接返回
	if newSize*8 == b.size {
		return nil
	}

	// 创建新数组并复制数据
	newBits := make([]byte, newSize)
	copy(newBits, b.bits) // 内置copy函数自动处理长度差异

	// 更新状态
	b.bits = newBits
	b.size = newSize * 8
	b.version++ // 修改版本号
	return nil
}

func Load(data []byte) (*Bitmap, error) {
	if len(data) == 0 {
		return NewBitmap(0)
	}
	if len(data) > maxSize {
		return nil, ErrInvalidSize
	}

	return &Bitmap{
		bits:    data,
		size:    len(data) * 8,
		version: 1,
	}, nil
}

func (b *Bitmap) Or(other *Bitmap) (*Bitmap, error) {
	if b.size != other.size {
		return nil, errors.New("bitmaps size mismatch")
	}

	b.mu.RLock()
	other.mu.RLock()
	defer b.mu.RUnlock()
	defer other.mu.RUnlock()

	result, _ := NewBitmap(len(b.bits))
	for i := range b.bits {
		result.bits[i] = b.bits[i] | other.bits[i]
	}
	return result, nil
}

/****************************** 辅助方法 ******************************/

// hash 优化的BKDR哈希函数
// 设计要点：
//   - 使用素数5381作为初始值（经验值，分布性好）
//   - hash*33 + c 的快速计算方式（通过移位优化）
//   - 最终结果取模确保在有效范围内
func (b *Bitmap) hash(id string) int {
	var hash uint32 = 5381 // 魔法素数初始值
	for _, c := range id {
		// 等价于 hash*33 + c，但通过移位优化性能
		hash = ((hash << 5) + hash) + uint32(c)
	}
	// 使用位操作代替取模（要求size是2的幂次）
	return int(hash & uint32(b.size-1))
}

/****************************** 二进制操作 ******************************/

// And 两个位图的按位与操作（线程安全）
// 注意：
//   - 要求两个位图大小相同
//   - 返回新位图不影响原数据
func (b *Bitmap) And(other *Bitmap) (*Bitmap, error) {
	// 大小校验
	if b.size != other.size {
		return nil, errors.New("bitmaps size mismatch")
	}

	// 双读锁保护
	b.mu.RLock()
	other.mu.RLock()
	defer b.mu.RUnlock()
	defer other.mu.RUnlock()

	// 创建结果位图
	result, _ := NewBitmap(len(b.bits))

	// 逐字节执行AND操作
	for i := range b.bits {
		result.bits[i] = b.bits[i] & other.bits[i]
	}
	return result, nil
}

// Export 导出位图数据（线程安全）
// 安全设计：
//   - 返回数据的深拷贝，避免外部修改影响内部状态
func (b *Bitmap) Export() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 创建新切片并复制数据
	data := make([]byte, len(b.bits))
	copy(data, b.bits)
	return data
}
