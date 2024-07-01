package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(str []byte) string {
	checksumHash := md5.Sum(str)
	return hex.EncodeToString(checksumHash[:])
}
