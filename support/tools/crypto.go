package tools

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

func Md5(in string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(in))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func Md5File(file string) string {
	md5 := md5.New()

	fi, err := os.Open(file)
	if err != nil {
		panic(err.Error())
	}
	io.Copy(md5, fi)

	return hex.EncodeToString(md5.Sum(nil))
}
