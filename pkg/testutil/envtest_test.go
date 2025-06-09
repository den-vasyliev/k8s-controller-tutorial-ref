package testutil

import "testing"

func TestInt32Ptr(t *testing.T) {
	v := int32(42)
	ptr := int32Ptr(v)
	if ptr == nil || *ptr != v {
		t.Errorf("int32Ptr(%d) = %v, want pointer to %d", v, ptr, v)
	}
}
