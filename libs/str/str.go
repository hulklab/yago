package str

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

// md5 加密字符串
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// sha1 加密字符串
func Sha1(s string) string {
	h := sha1.New()
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

func StrToKv(s, sep, pair string) map[string]string {
	m := make(map[string]string)

	// 先按 pair 分隔
	ss := strings.Split(s, pair)
	for _, v := range ss {
		// 再按 sep 分隔
		ms := strings.Split(v, sep)
		if len(ms) == 2 {
			m[ms[0]] = ms[1]
		}
	}

	return m
}

func KvToStr(m map[string]string, sep, pair string, order ...string) string {
	if len(m) == 0 {
		return ""
	}

	s := ""

	if len(order) > 0 {
		o := strings.ToLower(order[0])

		keys := make([]string, 0, len(m))
		for k, _ := range m {
			keys = append(keys, k)
		}

		if o == "asc" {
			sort.Strings(keys)
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(keys)))
		}

		for _, key := range keys {
			s += fmt.Sprintf("%s%s%s%s", key, sep, m[key], pair)
		}

	} else {
		for k, v := range m {
			s += fmt.Sprintf("%s%s%s%s", k, sep, v, pair)
		}
	}

	s = strings.TrimRight(s, pair)

	return s
}

func Split(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		cutSet := ", \n\r\t"
		for _, v := range cutSet {
			if r == rune(v) {
				return true
			}
		}
		return false
	})

}
