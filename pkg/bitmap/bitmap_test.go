package bitmap

import "testing"

// bitmap-单元测试
func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(100)
	b.Set("aaaa")
	b.Set("bbbb")
	b.Set("6666")
	b.Set("eeee")
	b.Set("8888")

	for _, bit := range b.bits {
		t.Logf("%b,%v", bit, bit)
	}
}
