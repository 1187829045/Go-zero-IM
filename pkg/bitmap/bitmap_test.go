package bitmap

import (
	"fmt"
	"testing"
)

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(5)

	b.Set("pppp")
	b.Set("222")
	b.Set("pppp")
	b.Set("ccc")
	b.Set("eee")
	b.Set("fff")
	for _, bit := range b.bits {
		fmt.Println(bit)
		t.Logf("%b, %v", bit, bit)
	}
}
