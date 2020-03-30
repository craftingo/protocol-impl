package protocol_impl

import "crypto/cipher"

// CFB8 represents a CFB8 implementation
type CFB8 struct {
	block              cipher.Block
	blockSize          int
	initialVector, tmp []byte
	decrypt            bool
}

// NewCFB8Decrypt returns a CFB8 decryption implementation
func NewCFB8Decrypt(block cipher.Block, initialVector []byte) *CFB8 {
	ivCopy := make([]byte, len(initialVector))
	copy(ivCopy, initialVector)
	return &CFB8{
		block:         block,
		blockSize:     block.BlockSize(),
		initialVector: ivCopy,
		tmp:           make([]byte, block.BlockSize()),
		decrypt:       true,
	}
}

// NewCFB8Encrypt returns a CFB8 encryption implementation
func NewCFB8Encrypt(block cipher.Block, initialVector []byte) *CFB8 {
	ivCopy := make([]byte, len(initialVector))
	copy(ivCopy, initialVector)
	return &CFB8{
		block:         block,
		blockSize:     block.BlockSize(),
		initialVector: ivCopy,
		tmp:           make([]byte, block.BlockSize()),
		decrypt:       false,
	}
}

// XORKeyStream performs an XOR key stream
func (cfb *CFB8) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		val := src[i]
		copy(cfb.tmp, cfb.initialVector)
		cfb.block.Encrypt(cfb.initialVector, cfb.initialVector)
		val = val ^ cfb.initialVector[0]

		copy(cfb.initialVector, cfb.tmp[1:])
		if cfb.decrypt {
			cfb.initialVector[15] = src[i]
		} else {
			cfb.initialVector[15] = val
		}

		dst[i] = val
	}
}
