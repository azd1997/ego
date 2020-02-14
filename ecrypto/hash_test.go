package ecrypto

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestHash_LargerThan(t *testing.T) {
	h1 := sha256.Sum256([]byte("你是最棒的"))
	h2 := sha256.Sum256([]byte("你是最菜的"))
	h3 := sha256.Sum256([]byte("你是最棒的"))
	H1 := Hash(h1[:])
	H2 := Hash(h2[:])
	H3 := Hash(h3[:])

	ret1 := H1.LargerThan(H2, "H1的地址", "H2的地址")
	ret2 := H1.LargerThan(H3, "H1的地址", "H3的地址")

	fmt.Println("ret1 = ", ret1, "; ret2 = ", ret2)
}
