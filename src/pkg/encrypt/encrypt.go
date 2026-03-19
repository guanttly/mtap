// Package encrypt 提供AES-256加解密和脱敏工具
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

// AESEncryptor AES-256-GCM加解密器
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor 创建加密器，key 必须为 32 字节（AES-256）
func NewAESEncryptor(key string) (*AESEncryptor, error) {
	k := []byte(key)
	if len(k) != 32 {
		return nil, errors.New("AES-256 key must be 32 bytes")
	}
	return &AESEncryptor{key: k}, nil
}

// Encrypt 加密明文，返回 base64 编码的密文
func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密 base64 编码的密文
func (e *AESEncryptor) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// MaskName 姓名脱敏（张三 → 张*，张三丰 → 张*丰，欧阳锋 → 欧**）
func MaskName(name string) string {
	runes := []rune(name)
	length := len(runes)
	if length == 0 {
		return ""
	}
	if length == 1 {
		return string(runes)
	}
	if length == 2 {
		return string(runes[0]) + "*"
	}
	// 3字及以上：首尾保留，中间用*
	masked := string(runes[0])
	for i := 1; i < length-1; i++ {
		masked += "*"
	}
	masked += string(runes[length-1])
	return masked
}

// MaskPhone 手机号脱敏：13812345678 → 138****5678
func MaskPhone(phone string) string {
	if utf8.RuneCountInString(phone) < 7 {
		return phone
	}
	return phone[:3] + strings.Repeat("*", len(phone)-7) + phone[len(phone)-4:]
}
