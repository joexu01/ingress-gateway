package public

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
)

//GenerateSaltPassword利用SHA-256和盐值字符串生成加盐口令
func GenerateSaltPassword(salt, password string) string {
	digest1:=sha256.New()
	digest1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", digest1.Sum(nil))

	digest2 := sha256.New()
	digest2.Write([]byte(str1+salt))
	return fmt.Sprintf("%x", digest2.Sum(nil))
}

//MD5 md5加密
func MD5(s string) string {
	h := md5.New()
	_, _ = io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Obj2Json(s interface{}) string {
	bts, _ := json.Marshal(s)
	return string(bts)
}
func InStringSlice(slice []string,str string) bool{
	for  _,item:=range slice{
		if str==item{
			return true
		}
	}
	return false
}