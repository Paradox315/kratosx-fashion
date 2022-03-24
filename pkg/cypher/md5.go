package cypher

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string, b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}
