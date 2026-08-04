package test

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func TestGetMd5(t *testing.T) {
	got := fmt.Sprintf("%x", md5.Sum([]byte("123456")))
	want := "e10adc3949ba59abbe56e057f20f883e"

	if got != want {
		t.Fatalf("GetMd5() = %q, want %q", got, want)
	}
}
