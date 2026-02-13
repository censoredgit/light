package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Sha1Hasher struct {
}

func (m Sha1Hasher) Sum(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}

func (m Sha1Hasher) BlockSize() int {
	return sha1.BlockSize
}
