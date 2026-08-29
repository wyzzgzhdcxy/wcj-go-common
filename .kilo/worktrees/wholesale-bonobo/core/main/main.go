package main

import (
	"log"
)

func main() {
	//encode1 := core.MyEncode{}
	//password := "mysecretpassword123456789012" // 加密密码（至少32字节）
	//now := time.Now()
	//sourceDir := "E:\\个人信息"               // 要打包的文件夹
	//encryptedFile := "E:\\test\\个人信息.tar" // 加密后的文件名
	//
	//// 1. 直接打包并加密，不创建临时文件
	//err := encode1.TarAndEncrypt(sourceDir, encryptedFile, password)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}

	//
	//encryptedFile := "E:\\test\\个人信息.tar" // 加密的文件
	//targetDir := "E:\\test\\个人信息"         // 解压到的目录
	//
	//// 1. 解密并解压文件
	////err := encode1.DecryptAndUntar(encryptedFile, targetDir, password)
	////if err != nil {
	////	fmt.Println("Error decrypting and untarring:", err)
	////	return
	////}
	//
	//fmt.Printf("Folder successfully tarred and encrypted without temporary files!%s\n", time.Since(now))

	//myMap := make(map[string]int)
	//
	//core.MapMergeInt(&myMap, "test", 2)
	//core.MapMergeInt(&myMap, "test", 2)
	//
	//fmt.Print(myMap)

	log.Printf("文件路径:%s,md5:%x", "1111", "222")
	log.Printf("文件路径:%s,md5:%x", "1111", "222")

}
