package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
)

var (
	encryptKey []byte
	key        []byte // AES密钥
)
var once sync.Once

func InitEncryptKey(key string) {
	once.Do(func() {
		// 确保密钥长度为32字节(AES-256)
		encryptKey = make([]byte, 32)
		copy(encryptKey, []byte(key))
	})
}

// Encrypt 加密字符串
func Encrypt(plaintext string) string {
	if plaintext == "" || len(encryptKey) == 0 {
		return plaintext
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext // 加密失败时返回原文
	}

	// 创建gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext
	}

	// 创建nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// base64编码
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt 解密字符串
func Decrypt(ciphertext string) string {
	if ciphertext == "" {
		return ""
	}

	// base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return ciphertext
	}

	// 创建cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertext
	}

	// 创建gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return ciphertext
	}

	// 提取nonce
	if len(data) < gcm.NonceSize() {
		return ciphertext
	}
	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return ciphertext
	}

	return string(plaintext)
}
