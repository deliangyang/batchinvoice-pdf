package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	// 凭据文件名
	credentialFileName = ".credentials"
	// 加密盐值
	encryptionSalt = "batchinvoice-pdf-2026-secure-salt"
)

// CredentialManager 凭据管理器
type CredentialManager struct {
	dataDir string
}

// NewCredentialManager 创建凭据管理器
func NewCredentialManager() *CredentialManager {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}
	return &CredentialManager{
		dataDir: dir,
	}
}

// SavePassword 保存加密的密码
func (cm *CredentialManager) SavePassword(email, password string) error {
	if email == "" || password == "" {
		return fmt.Errorf("email and password cannot be empty")
	}

	// 生成加密密钥（基于邮箱和盐值）
	key := cm.generateKey(email)

	// 加密密码
	encryptedPassword, err := cm.encrypt(password, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// 保存到文件
	credentialFile := filepath.Join(cm.dataDir, credentialFileName)
	// 格式：email|encrypted_password
	data := fmt.Sprintf("%s|%s", email, encryptedPassword)

	err = os.WriteFile(credentialFile, []byte(data), 0600)
	if err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

// LoadPassword 加载并解密密码
func (cm *CredentialManager) LoadPassword(email string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	credentialFile := filepath.Join(cm.dataDir, credentialFileName)

	// 读取文件
	data, err := os.ReadFile(credentialFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // 文件不存在，返回空密码
		}
		return "", fmt.Errorf("failed to read credentials: %w", err)
	}

	// 解析数据
	content := string(data)
	var savedEmail, encryptedPassword string
	n, err := fmt.Sscanf(content, "%s|%s", &savedEmail, &encryptedPassword)
	if err != nil || n != 2 {
		return "", fmt.Errorf("invalid credential file format")
	}

	// 验证邮箱是否匹配
	if savedEmail != email {
		return "", nil // 邮箱不匹配，返回空密码
	}

	// 生成解密密钥
	key := cm.generateKey(email)

	// 解密密码
	password, err := cm.decrypt(encryptedPassword, key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	return password, nil
}

// DeletePassword 删除保存的密码
func (cm *CredentialManager) DeletePassword() error {
	credentialFile := filepath.Join(cm.dataDir, credentialFileName)
	err := os.Remove(credentialFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete credentials: %w", err)
	}
	return nil
}

// generateKey 生成加密密钥
func (cm *CredentialManager) generateKey(email string) []byte {
	// 使用 SHA256 生成 32 字节的密钥（适合 AES-256）
	hash := sha256.Sum256([]byte(email + encryptionSalt))
	return hash[:]
}

// encrypt 使用 AES-GCM 加密数据
func (cm *CredentialManager) encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt 使用 AES-GCM 解密数据
func (cm *CredentialManager) decrypt(ciphertext string, key []byte) (string, error) {
	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
