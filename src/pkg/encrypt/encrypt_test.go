package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testKey = "01234567890123456789012345678901" // 32 bytes

func TestNewAESEncryptor_ValidKey(t *testing.T) {
	enc, err := NewAESEncryptor(testKey)
	require.NoError(t, err)
	assert.NotNil(t, enc)
}

func TestNewAESEncryptor_InvalidKey(t *testing.T) {
	_, err := NewAESEncryptor("too-short")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "32 bytes")
}

func TestEncryptDecrypt(t *testing.T) {
	enc, _ := NewAESEncryptor(testKey)

	tests := []string{
		"张三",
		"13812345678",
		"Hello, 世界!",
		"",
		"a",
		"这是一段很长的中文文本，用来测试加解密功能是否正常工作",
	}

	for _, plaintext := range tests {
		t.Run(plaintext, func(t *testing.T) {
			ciphertext, err := enc.Encrypt(plaintext)
			require.NoError(t, err)
			assert.NotEqual(t, plaintext, ciphertext)

			decrypted, err := enc.Decrypt(ciphertext)
			require.NoError(t, err)
			assert.Equal(t, plaintext, decrypted)
		})
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	enc, _ := NewAESEncryptor(testKey)
	c1, _ := enc.Encrypt("same input")
	c2, _ := enc.Encrypt("same input")
	// GCM uses random nonce, so ciphertexts differ
	assert.NotEqual(t, c1, c2)
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	enc, _ := NewAESEncryptor(testKey)
	_, err := enc.Decrypt("not-base64!!!")
	assert.Error(t, err)
}

func TestDecrypt_TamperedData(t *testing.T) {
	enc, _ := NewAESEncryptor(testKey)
	ciphertext, _ := enc.Encrypt("secret")
	// tamper with ciphertext
	tampered := ciphertext[:len(ciphertext)-2] + "XX"
	_, err := enc.Decrypt(tampered)
	assert.Error(t, err)
}

func TestDecrypt_WrongKey(t *testing.T) {
	enc1, _ := NewAESEncryptor("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	enc2, _ := NewAESEncryptor("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	ciphertext, _ := enc1.Encrypt("secret")
	_, err := enc2.Decrypt(ciphertext)
	assert.Error(t, err)
}

func TestMaskName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"张", "张"},
		{"张三", "张*"},
		{"张三丰", "张*丰"},
		{"欧阳修改", "欧**改"},
		{"司马相如说", "司***说"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, MaskName(tt.input))
		})
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"13812345678", "138****5678"},
		{"123", "123"},       // too short, no masking
		{"12345678901", "123****8901"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, MaskPhone(tt.input))
		})
	}
}
