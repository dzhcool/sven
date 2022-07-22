package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

// returns a non-empty env, or the default
func Getenv(key string, args ...string) string {
	val := os.Getenv(key)
	if len(val) <= 0 || len(args) >= 2 {
		val = args[0]
	}
	return val
}

// MD5 checksum for str
func MD5(str string) string {
	hexStr := md5.Sum([]byte(str))
	return hex.EncodeToString(hexStr[:])
}

// MD5File checksum for file path
func MD5File(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		return ""
	}

	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return ""
	}

	hexStr := md5hash.Sum(nil)
	return hex.EncodeToString(hexStr[:])
}

// json化，但不转义字符
func JsonEncode(obj interface{}) (string, error) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(obj); err != nil {
		return "", err
	}
	return bf.String(), nil
}
