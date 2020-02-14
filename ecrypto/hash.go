package ecrypto

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

/*********************************************************************************************************************
                                                    Hash相关
*********************************************************************************************************************/

// Hash 32B哈希。如果要修改哈希算法，只需在这里重新定义哈希的具体类型即可
// 使用[32]byte ，使用起来太不方便。
type Hash []byte

// String 打印方法
func (h Hash) String() string {
	// 直接转为string会出现乱码，所以还是转为16进制打印
	return fmt.Sprintf("%x", string(h))
}

// LargerThan 判断哈希A是否大于哈希B。
// 如果遇到相等，则不断加salt（ah+asalt, bh+bsalt）再哈希再比较，直至大于或小于
// 例如在Ecare项目中，asalt,bsalt使用双方的通信地址字符串
func (h Hash) LargerThan(bh Hash, asalt, bsalt string) bool {
	hcopy := h
	var ahs, bhs string
	var h32, bh32 [32]byte
	for {
		ahs, bhs = string(hcopy), string(bh)
		if ahs > bhs {
			return true
		} else if ahs < bhs {
			return false
		} else {	// ahs==bhs
			h32 = sha256.Sum256(append(hcopy, []byte(asalt)...))
			bh32 = sha256.Sum256(append(bh, []byte(bsalt)...))
			hcopy, bh = h32[:], bh32[:]
		}
	}
}

// BytesToHash 将长度为32的字节切片转换为Hash，若返回Hash{}，说明有错
func BytesToHash(data []byte) [32]byte {
	var res [32]byte
	if len(data) != cap(res) {
		return [32]byte{}	// 若返回Hash{}，说明有错
	}

	for i := 0; i < cap(res); i++ {
		res[i] = data[i]
	}

	return res
}

// RandomHash 生成随机的Hash。只是用来作一些测试
func RandomHash() Hash {
	res := make([]byte, 32)
	for i:=0; i < 32; i++ {
		res[i] = byte(uint(rand.Intn(256)))
	}
	return res
}

// ZeroHASH 全局零哈希变量
var ZeroHASH = ZeroHash()

// ZeroHash 生成全0哈希
func ZeroHash() (zero Hash) {
	zero = make([]byte, 32)
	for i:=0; i < 32; i++ {
		zero[i] = byte(0)
	}
	return zero
}
