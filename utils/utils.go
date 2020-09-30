package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"hash/crc32"
)

// MD5 生成md5
func MD5(str string) string {
	c := md5.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

// SHA1 生成sha1
func SHA1(str string) string {
	c := sha1.New()
	c.Write([]byte(str))
	return hex.EncodeToString(c.Sum(nil))
}

// CRC32 生成CRC32
func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// GetSession 获取session
func GetSession(c *gin.Context, key string) interface{} {
	session := sessions.Default(c)
	v := session.Get(key)

	return v
}
