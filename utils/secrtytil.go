package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz!@#$%^ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := time.Now().UnixNano()
	randSrc := rand.New(rand.NewSource(seed))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randSrc.Intn(len(charset))]
	}
	return string(b)
}

func GenerateMD5Hash(input string) string {
	combinedString := input

	// 创建 MD5 哈希对象
	hasher := md5.New()

	// 将组合字符串转换为字节数组并计算 MD5 哈希值
	hasher.Write([]byte(combinedString))
	hashBytes := hasher.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func Md5Crypt(str string, salt ...interface{}) (CryptStr string) {
	if l := len(salt); l > 0 {
		slice := make([]string, l+1)
		str = fmt.Sprintf(str+strings.Join(slice, "%v"), salt...)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
