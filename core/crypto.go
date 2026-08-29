package core

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"os"

	"github.com/tjfoc/gmsm/sm2"
)

// 编码、哈希、加密相关的工具方法。
// 压缩/归档相关方法见 archive.go。

// Base64Encode 编码base64
func Base64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func Base64EncodeStr(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

// Base64Decode unused
// Base64Decode 解码base64
func Base64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func Base64DecodeStr(src string) string {
	byteContent, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		log.Printf("Base64DecodeString error,%v\n", err)
	}
	return string(byteContent)
}

func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func Md5File(filePath string) string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer Close(file)

	// 计算 MD5
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}

	// 获取 MD5 哈希值（16进制字符串）
	md5Sum := hash.Sum(nil)
	log.Printf("文件路径:%s,md5:%x", filePath, md5Sum)
	return hex.EncodeToString(md5Sum)
}

func Md5Byte(plainByte []byte) string {
	hash := md5.Sum(plainByte)
	return hex.EncodeToString(hash[:])
}

func TestSm2Encrypt() {
	random := rand.Reader //If there is no external trusted random source,please use rand.Reader to instead of it.
	//生成私钥
	privateKey, e := sm2.GenerateKey(random)
	if e != nil {
		log.Printf("sm2 encrypt failed！")
	}
	//从私钥中获取公钥
	pubkey := &privateKey.PublicKey
	msg := []byte("i am   wek && 政府第三。")
	//用公钥加密msg
	encryptBytes, err := pubkey.EncryptAsn1(msg, random)
	if err != nil {
		log.Printf("使用私钥加密失败！")
	}
	log.Printf("the encrypt msg  =  ", hex.EncodeToString(encryptBytes))
	//用私钥解密msg
	decrypt, i2 := privateKey.DecryptAsn1(encryptBytes)
	if i2 != nil {
		log.Printf("使用私钥解密失败！")
	}
	log.Printf("the msg  = ", string(decrypt))
}
