package osm

import (
	"crypto/md5"
	"fmt"
)

func GetMD5Hash(data string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		fmt.Println("Error generating MD5 hash:", err)
		return ""
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
