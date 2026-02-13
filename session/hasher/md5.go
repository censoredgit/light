package hasher

import (
	"crypto/md5"
	"fmt"
)

type Md5Hasher struct {
}

func (m Md5Hasher) Sum(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

func (m Md5Hasher) BlockSize() int {
	return md5.BlockSize
}
