package osm

import (
	"fmt"
	"hash/fnv"
)

func GetHash(data string) string {
	h := fnv.New64a()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum64())
}
