package core

import (
	"math/rand"
	"strconv"
	"time"
)

func RandomInt(count int) string {
	var r string
	for i := 0; i < count; i++ {
		tmp := rand.Intn(10)
		r = r + strconv.Itoa(tmp)
	}
	return r
}

func GetRandomString(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
