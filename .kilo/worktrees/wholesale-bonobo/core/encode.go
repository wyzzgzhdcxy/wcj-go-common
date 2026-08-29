package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"log"

	"github.com/tjfoc/gmsm/sm2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//编码,加密相关工具类

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

// DecompressZip src: 需要解压的zip文件
func DecompressZip(src string) error {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	dir := filepath.Dir(src)
	defer Close(archive)
	for _, f := range archive.File {
		filePath := filepath.Join(dir, getFilename(f))
		if f.FileInfo().IsDir() {
			MkDirALl0755(filePath)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to make directory (%v)", err)
		}
		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file (%v)", err)
		}
		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip (%v)", err)
		}
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("failed to copy file in zip (%v)", err)
		}
		Close(dstFile)
		Close(&fileInArchive)
	}
	return nil
}

// 获取压缩文件文件名，解决中文乱码问题
func getFilename(f *zip.File) string {
	if f.Flags != 0 {
		//如果标志为是 1 << 11也就是 2048  则是utf-8编码
		return f.Name
	}
	//如果标致位是0  则是默认的本地编码   默认为gbk
	i := bytes.NewReader([]byte(f.Name))
	decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
	content, _ := io.ReadAll(decoder)
	return string(content)
}

// dir: 需要打包的本地文件夹路径
// dst: 保存压缩包的本地路径
func CompressDir(dir string, dst string) error {
	zipFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer Close(zipFile)
	archive := zip.NewWriter(zipFile)
	defer Close(archive)
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir {
			return nil
		}
		info, _ := d.Info()
		h, _ := zip.FileInfoHeader(info)
		h.Name = strings.TrimPrefix(path, dir+"/")
		if info.IsDir() {
			h.Name += "/"
		} else {
			h.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(h)
		if !info.IsDir() {
			srcFile, _ := os.Open(path)
			defer Close(srcFile)
			_, _ = io.Copy(writer, srcFile)
		}
		return nil
	})
}

// GzipEncode unused
func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		Close(writer)
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer Close(reader)
	return io.ReadAll(reader)
}

// TarAndEncrypt 流式打包并加密
func TarAndEncrypt(sourceDir, outputFile, password string) error {
	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer Close(outFile)

	// 生成密钥（确保密码至少32字节）
	key := make([]byte, 32) // AES-256
	copy(key, password)
	for i := len(password); i < len(key); i++ {
		key[i] = 0
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// 生成随机的 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	// 写入 IV 到输出文件
	if _, err := outFile.Write(iv); err != nil {
		return err
	}

	// 创建加密流
	stream := cipher.NewCFBEncrypter(block, iv)

	// 创建管道
	pr, pw := io.Pipe()
	defer Close(pr)

	// 启动goroutine进行打包
	go func() {
		gw := gzip.NewWriter(pw)
		tw := tar.NewWriter(gw)

		err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 创建 tar 头信息
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			// 确保路径是相对于源目录的
			relPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}
			header.Name = relPath

			// 写入头信息
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			// 如果是普通文件，写入文件内容
			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer Close(file)
				if _, err := io.Copy(tw, file); err != nil {
					return err
				}
			}
			return nil
		})

		// 确保所有写入器都关闭
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		if err := tw.Close(); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		if err := gw.Close(); err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		_ = pw.Close()
	}()

	// 加密并写入文件
	buf := make([]byte, 4096)
	for {
		n, err := pr.Read(buf)
		if n > 0 {
			// 正确使用 XORKeyStream
			dst := make([]byte, n)
			stream.XORKeyStream(dst, buf[:n])
			if _, err := outFile.Write(dst); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// DecryptAndUntar 解密并解压文件
func DecryptAndUntar(encryptedFile, targetDir, password string) error {
	// 打开加密文件
	inFile, err := os.Open(encryptedFile)
	if err != nil {
		return err
	}
	defer Close(inFile)

	// 生成密钥（确保与加密时一致）
	key := make([]byte, 32) // AES-256
	copy(key, password)
	for i := len(password); i < len(key); i++ {
		key[i] = 0
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// 读取 IV（前16字节）
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(inFile, iv); err != nil {
		return err
	}

	// 创建解密流
	stream := cipher.NewCFBDecrypter(block, iv)

	// 创建管道
	pr, pw := io.Pipe()
	defer Close(pr)

	// 启动goroutine进行解密
	go func() {
		defer Close(pw)
		buf := make([]byte, 4096)
		for {
			n, err := inFile.Read(buf)
			if n > 0 {
				// 解密数据
				dst := make([]byte, n)
				stream.XORKeyStream(dst, buf[:n])
				if _, err := pw.Write(dst); err != nil {
					_ = pw.CloseWithError(err)
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				_ = pw.CloseWithError(err)
				return
			}
		}
	}()

	// 解压tar.gz
	gr, err := gzip.NewReader(pr)
	if err != nil {
		return err
	}
	defer Close(gr)
	tr := tar.NewReader(gr)
	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// 解压文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 文件结束
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, header.Name)

		// 根据文件类型处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// 创建文件
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				Close(outFile)
				return err
			}
			Close(outFile)
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

// Tar 打包目录到tar文件 要打包的目录路径
// 目标tar文件路径(.tar或.tar.gz) 是否使用gzip压缩
func Tar(srcDir, destFile string, useGzip bool) error {
	// 检查源目录是否存在
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在: %s", srcDir)
	}
	// 创建目标文件
	file, err := os.Create(destFile)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer Close(file)
	var tw *tar.Writer
	// 根据是否使用gzip创建不同的writer
	if useGzip {
		gw := gzip.NewWriter(file)
		defer Close(gw)
		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(file)
	}
	defer Close(tw)
	// 遍历目录并添加到tar
	return filepath.Walk(srcDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 创建tar头信息
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}
		// 调整头信息中的名称，使其相对于源目录
		relPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath) // 确保使用正斜杠
		fmt.Printf("打包: %s\n", header.Name)
		// 写入头信息
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// 如果是普通文件，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer Close(file)
			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UnTar 解压tar文件到指定目录 要解压的tar文件路径 目标目录路径
func UnTar(srcFile, destDir string) error {
	// 检查源文件是否存在
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %s", srcFile)
	}

	// 创建目标目录(如果不存在)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 打开tar文件
	file, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("打开tar文件失败: %v", err)
	}
	defer Close(file)

	var tr *tar.Reader

	// 检查是否是gzip压缩
	if strings.HasSuffix(srcFile, ".gz") || strings.HasSuffix(srcFile, ".tgz") {
		gzr, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("创建gzip reader失败: %v", err)
		}
		defer Close(gzr)
		tr = tar.NewReader(gzr)
	} else {
		tr = tar.NewReader(file)
	}

	// 解压每个文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 归档结束
		}
		if err != nil {
			return fmt.Errorf("读取tar头信息失败: %v", err)
		}
		target := filepath.Join(destDir, header.Name)
		fmt.Printf("解压: %s\n", header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
		case tar.TypeReg:
			// 确保父目录存在
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("创建父目录失败: %v", err)
			}
			// 创建文件
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("创建文件失败: %v", err)
			}
			// 写入文件内容
			if _, err := io.Copy(outFile, tr); err != nil {
				Close(outFile)
				return fmt.Errorf("写入文件内容失败: %v", err)
			}
			Close(outFile)
		default:
			return fmt.Errorf("不支持的文件类型 %v 在 %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
