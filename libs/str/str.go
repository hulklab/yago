package str

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

// md5 加密字符串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var node, _ = snowflake.NewNode(int64(rand.Intn(1024)))

// 全局唯一 id
func UniqueId() string {
	return Md5(node.Generate().String())
}

// 全局唯一 短id
func UniqueIdShort() string {
	return node.Generate().Base58()
}

func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 生成时间的唯一 20 长度的 id
func GenJid() string {
	date := time.Now().Format("060102150405")
	micro := time.Now().Nanosecond() / 10
	jid := fmt.Sprintf("%s%08d", date, micro)
	time.Sleep(time.Duration(10) * time.Nanosecond)
	return jid
}

// 首字母大写
func Ucfirst(s string) string {
	if len(s) < 1 {
		return ""
	}
	strSlice := []rune(s)
	if strSlice[0] >= 97 && strSlice[0] <= 122 {
		strSlice[0] -= 32
	}
	return string(strSlice)
}
