package bitmap

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// TestNewBitmap 测试创建新位图
func TestNewBitmap(t *testing.T) {
	t.Run("正常创建", func(t *testing.T) {
		b, err := NewBitmap(100)
		if err != nil {
			t.Fatalf("创建位图失败: %v", err)
		}
		if len(b.bits) != 100 {
			t.Errorf("期望大小100字节，实际得到%d字节", len(b.bits))
		}
	})

	t.Run("默认大小", func(t *testing.T) {
		b, err := NewBitmap(0)
		if err != nil {
			t.Fatal(err)
		}
		if len(b.bits) != defaultSize {
			t.Errorf("期望默认大小%d字节，实际得到%d字节", defaultSize, len(b.bits))
		}
	})

	t.Run("非法参数", func(t *testing.T) {
		_, err := NewBitmap(-1)
		if err != ErrInvalidSize {
			t.Errorf("期望错误%v，实际得到%v", ErrInvalidSize, err)
		}

		_, err = NewBitmap(maxSize + 1)
		if err != ErrInvalidSize {
			t.Errorf("期望错误%v，实际得到%v", ErrInvalidSize, err)
		}
	})
}

// TestSetAndIsSet 测试设置和检查bit
func TestSetAndIsSet(t *testing.T) {
	b, _ := NewBitmap(32) // 256 bits

	t.Run("基本设置检查", func(t *testing.T) {
		// 测试边界值
		testPositions := []int{0, 1, 63, 64, 127, 255}
		for _, pos := range testPositions {
			if err := b.SetBit(pos); err != nil {
				t.Fatalf("设置位置%d失败: %v", pos, err)
			}

			set, err := b.IsBitSet(pos)
			if err != nil {
				t.Fatal(err)
			}
			if !set {
				t.Errorf("位置%d应该被设置但未被检测到", pos)
			}
		}
	})

	t.Run("字符串ID设置检查", func(t *testing.T) {
		ids := []string{"user1", "item42", "session_abc123"}
		for _, id := range ids {
			if err := b.Set(id); err != nil {
				t.Fatal(err)
			}

			set, err := b.IsSet(id)
			if err != nil {
				t.Fatal(err)
			}
			if !set {
				t.Errorf("ID %s应该被设置但未被检测到", id)
			}
		}
	})

	t.Run("越界检查", func(t *testing.T) {
		if err := b.SetBit(256); err != ErrIndexOutOfRange {
			t.Errorf("期望错误%v，实际得到%v", ErrIndexOutOfRange, err)
		}

		if err := b.SetBit(-1); err != ErrIndexOutOfRange {
			t.Errorf("期望错误%v，实际得到%v", ErrIndexOutOfRange, err)
		}
	})
}

// TestClearAndReset 测试清除和重置
func TestClearAndReset(t *testing.T) {
	b, _ := NewBitmap(10)

	t.Run("清除单个bit", func(t *testing.T) {
		pos := 42
		b.SetBit(pos) // 先设置

		if err := b.ClearBit(pos); err != nil {
			t.Fatal(err)
		}

		set, _ := b.IsBitSet(pos)
		if set {
			t.Errorf("位置%d应该被清除但仍被检测为设置", pos)
		}
	})

	t.Run("重置整个位图", func(t *testing.T) {
		// 设置多个bit
		for i := 0; i < 50; i++ {
			b.SetBit(i)
		}

		b.Reset()

		for i := 0; i < 50; i++ {
			set, _ := b.IsBitSet(i)
			if set {
				t.Errorf("位置%d在Reset后应该被清除", i)
			}
		}
	})
}

// TestCount 测试统计功能
func TestCount(t *testing.T) {
	b, _ := NewBitmap(10) // 80 bits

	// 设置10个随机的bit
	expected := 10
	for i := 0; i < expected; i++ {
		pos := rand.Intn(80)
		b.SetBit(pos)
	}

	// 可能有重复位置，所以实际设置数 <= expected
	count := b.Count()
	if count > expected {
		t.Errorf("期望最多%d个bit被设置，实际得到%d", expected, count)
	}
	if count == 0 {
		t.Error("至少应该有一个bit被设置")
	}
}

// TestResize 测试调整大小
func TestResize(t *testing.T) {
	b, _ := NewBitmap(5) // 40 bits

	// 设置一些bit
	b.SetBit(0)
	b.SetBit(39)

	t.Run("扩大尺寸", func(t *testing.T) {
		if err := b.Resize(10); err != nil {
			t.Fatal(err)
		}

		// 检查原有bit是否保留
		if set, _ := b.IsBitSet(0); !set {
			t.Error("扩大后位置0应该保留")
		}
		if set, _ := b.IsBitSet(39); !set {
			t.Error("扩大后位置39应该保留")
		}
	})

	t.Run("缩小尺寸", func(t *testing.T) {
		if err := b.Resize(3); err != nil {
			t.Fatal(err)
		}

		// 检查超出范围的bit是否被丢弃
		if _, err := b.IsBitSet(39); err != ErrIndexOutOfRange {
			t.Error("缩小后超出范围的访问应该报错")
		}
	})

	t.Run("非法尺寸", func(t *testing.T) {
		if err := b.Resize(-1); err != ErrInvalidSize {
			t.Errorf("期望错误%v，实际得到%v", ErrInvalidSize, err)
		}
	})
}

// TestExportAndLoad 测试导出和加载
func TestExportAndLoad(t *testing.T) {
	b1, _ := NewBitmap(10)

	// 设置一些bit
	for i := 0; i < 20; i++ {
		b1.SetBit(rand.Intn(80))
	}

	// 导出数据
	data := b1.Export()

	// 从数据加载新位图
	b2, err := Load(data)
	if err != nil {
		t.Fatal(err)
	}

	// 比较所有bit是否一致
	for i := 0; i < 80; i++ {
		set1, _ := b1.IsBitSet(i)
		set2, _ := b2.IsBitSet(i)
		if set1 != set2 {
			t.Errorf("位置%d在两个位图中不一致", i)
		}
	}
}

// TestConcurrentAccess 测试并发安全性
func TestConcurrentAccess(t *testing.T) {
	b, _ := NewBitmap(100)
	var wg sync.WaitGroup

	// 并发设置
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			b.SetBit(pos % 800) // 800 bits范围内
			b.IsBitSet(pos % 800)
		}(i)
	}

	// 并发计数
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = b.Count()
		}()
	}

	wg.Wait()

	// 最终检查没有panic即认为测试通过
}

// BenchmarkBitmap 性能基准测试
func BenchmarkBitmap(b *testing.B) {
	// 测试不同操作的性能
	b.Run("SetBit", func(b *testing.B) {
		bm, _ := NewBitmap(1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = bm.SetBit(i % 8000)
		}
	})

	b.Run("IsBitSet", func(b *testing.B) {
		bm, _ := NewBitmap(1000)
		for i := 0; i < 1000; i++ {
			bm.SetBit(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = bm.IsBitSet(i % 8000)
		}
	})

	b.Run("Count", func(b *testing.B) {
		bm, _ := NewBitmap(1000)
		for i := 0; i < 500; i++ {
			bm.SetBit(rand.Intn(8000))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = bm.Count()
		}
	})

	b.Run("Concurrent", func(b *testing.B) {
		bm, _ := NewBitmap(10000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				pos := rand.Intn(80000)
				if rand.Float32() < 0.5 {
					_ = bm.SetBit(pos)
				} else {
					_, _ = bm.IsBitSet(pos)
				}
			}
		})
	})
}

// TestBitOperations 测试位图操作
func TestBitOperations(t *testing.T) {
	b1, _ := NewBitmap(10)
	b2, _ := NewBitmap(10)

	// 设置b1的偶数位
	for i := 0; i < 80; i += 2 {
		b1.SetBit(i)
	}

	// 设置b2的3的倍数位
	for i := 0; i < 80; i += 3 {
		b2.SetBit(i)
	}

	t.Run("AND操作", func(t *testing.T) {
		result, err := b1.And(b2)
		if err != nil {
			t.Fatal(err)
		}

		// 检查结果应该是6的倍数位
		for i := 0; i < 80; i++ {
			set, _ := result.IsBitSet(i)
			shouldSet := i%6 == 0
			if set != shouldSet {
				t.Errorf("位置%d应该为%v，实际为%v", i, shouldSet, set)
			}
		}
	})

	t.Run("OR操作", func(t *testing.T) {
		result, err := b1.Or(b2)
		if err != nil {
			t.Fatal(err)
		}

		// 检查结果应该是2或3的倍数位
		for i := 0; i < 80; i++ {
			set, _ := result.IsBitSet(i)
			shouldSet := i%2 == 0 || i%3 == 0
			if set != shouldSet {
				t.Errorf("位置%d应该为%v，实际为%v", i, shouldSet, set)
			}
		}
	})

	t.Run("大小不匹配", func(t *testing.T) {
		b3, _ := NewBitmap(5)
		_, err := b1.And(b3)
		if err == nil {
			t.Error("期望大小不匹配错误但未收到")
		}
	})
}

// TestHashDistribution 测试哈希分布
func TestHashDistribution(t *testing.T) {
	b, _ := NewBitmap(100) // 800 bits
	bucketSize := 10
	buckets := make([]int, bucketSize)

	// 测试1000个不同的字符串ID
	for i := 0; i < 1000; i++ {
		id := "item_" + strconv.Itoa(i)
		pos := b.hash(id)
		bucket := pos % bucketSize
		buckets[bucket]++
	}

	// 检查分布是否均匀（每个桶应该有大约100个元素）
	expected := 100
	tolerance := 30 // 允许±30的偏差
	for i, count := range buckets {
		if count < expected-tolerance || count > expected+tolerance {
			t.Errorf("桶%d分布不均匀，期望约%d，实际%d", i, expected, count)
		}
	}
}
