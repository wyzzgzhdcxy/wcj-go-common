package core

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"hash/crc32"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

// MkDirALl0755 0700：表示所有者有读、写和执行权限，组和其他人没有任何权限。
// 0755：表示所有者有读、写和执行权限，组和其他人有读和执行权限。
// 0644：表示所有者有读和写权限，组和其他人有读权限。
func MkDirALl0755(filepathList ...string) {
	MkDirAll(0755, filepathList...)
}
func MkDirALl0777(filepathList ...string) {
	MkDirAll(0777, filepathList...)
}

func MkDirAll(perm os.FileMode, filepathList ...string) {
	for _, fp := range filepathList {
		err := os.MkdirAll(fp, perm)
		if err != nil {
			log.Printf("创建文件失败:", fp, "，错误信息", err)
		}
	}
	fmt.Printf("目录创建成功:%v\n", filepathList)
}

func MergeFile(dirPath string) {
	outputFile, err := os.Create(dirPath + ".apk")
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer Close(outputFile)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer Close(file)
			_, err = io.Copy(outputFile, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func SplitFile(filePath string, chunkSize int64) {
	srcFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("%v", err)
	}
	defer Close(srcFile)
	filename := filepath.Base(filePath)
	dir := filepath.Dir(filePath)
	for i := 0; ; i++ {
		chunkFile, err := os.Create(dir + "/1111/" + filename + "." + strconv.Itoa(i))
		if err != nil {
			log.Printf("%v", err)
		}
		defer Close(chunkFile)
		_, err = io.CopyN(chunkFile, srcFile, chunkSize)
		if err != nil && err != io.EOF {
			log.Printf("%v", err)
			return
		}
		// If reached EOF, and last chunk was less than chunkSize, adjust to just the size of the last chunk
		if err == io.EOF {
			//chunkSize, _ = srcFile.Seek(0, io.SeekCurrent)
			return
		}
	}
}

// FindFirstChineseRuneIndex 查找第一个中文字符的索引位置
func FindFirstChineseRuneIndex(name string) int {
	for i, runeValue := range name {
		if unicode.Is(unicode.Scripts["Han"], runeValue) {
			return i
		}
	}
	return -1 // 如果没有找到中文字符，返回-1
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// IsDirEmpty 判断目录是否为空
func IsDirEmpty(dirPath string) bool {
	// 打开目录
	dir, err := os.Open(dirPath)
	if err != nil {
		return true
	}
	defer dir.Close()

	// 读取目录中的第一个条目
	// 如果目录为空，ReadDir 会返回空切片，且没有错误
	entries, err := dir.ReadDir(1) // 只读取一个条目，用于判断是否为空
	if err != nil {
		return true
	}

	// 如果没有条目，则目录为空
	return len(entries) == 0
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FormatSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"} // 单位
	i := 0                                               // 计数器
	for size >= 1024 && i < len(units)-1 {
		size /= 1024
		i++
	}
	return fmt.Sprintf("%d%s", size, units[i])
}

func IsDir(path string) bool {
	// 使用os.Stat函数获取文件或目录的信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("无法获取文件或目录的信息:%v", err)
		panic("无法获取文件或目录的信息")
	}
	// 判断文件或目录的类型
	return fileInfo.IsDir()
}

// DeleteDir 删除指定目录及其所有内容（递归）
func DeleteDir(dir string) {
	DeleteFilesInDir(dir)
	// 删除空目录
	err := os.Remove(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// DeleteFilesInDir 删除指定目录及其所有内容（递归）
func DeleteFilesInDir(dir string) {
	// 打开目录
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer Close(d)

	// 读取目录内容
	entries, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 遍历目录内容
	for _, entry := range entries {
		// 构建完整路径
		path := filepath.Join(dir, entry.Name())
		// 如果是目录，则递归删除
		if entry.IsDir() {
			DeleteDir(path)
		} else {
			// 如果是文件，则删除
			err = os.Remove(path)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}

func DeleteFile(filePath string) {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("文件不存在: %s\n", filePath)
			return
		}
		fmt.Printf("无法获取文件信息: %s\n", filePath)
		return
	}
	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("无法删除文件: %s\n", filePath)
		return
	}
	fmt.Printf("文件已成功删除: %s\n", filePath)
}

func ReadFileToStr(filepath string) string {
	return string(ReadFileToByte(filepath))
}

func ReadFileToByte(filepath string) []byte {
	// 读取hosts文件内容
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading hosts file: %v", err)
		return []byte("")
	}
	return data
}

// WriteStrToFile 写文件
func WriteStrToFile(filepath string, text string) {
	err := os.WriteFile(filepath, []byte(text), 0644)
	if err != nil {
		fmt.Printf("写入文件失败,filepath:%s,err:%v\n", filepath, err)
	}
}

func WriteByteToFile(filepath string, content []byte) {
	err := os.WriteFile(filepath, content, 0644)
	if err != nil {
		fmt.Printf("写入文件失败,filepath:%s,err:%v\n", filepath, err)
	}
}

// AppendStrToFile 添加字符串到文件末尾
func AppendStrToFile(filePath string, appendStr string) {
	// 打开文件，如果文件不存在则创建新文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("无法打开文件: %s\n", err)
		return
	}
	defer Close(file)
	// 向文件追加字节数组
	_, err = file.Write([]byte(appendStr))
	if err != nil {
		fmt.Printf("无法向文件追加字符: %s\n", err)
		return
	}
}

// WriteStrToFileGBK 写文件
func WriteStrToFileGBK(filepath string, text string) {
	data := []byte(text)
	// 创建一个新的字节缓冲区
	var buf bytes.Buffer

	// 创建一个 GBK 编码器
	writer := transform.NewWriter(&buf, simplifiedchinese.GBK.NewEncoder())

	// 将数据写入编码器
	_, err := writer.Write(data)
	if err != nil {
		fmt.Println("写入编码器时出错:", err)
		return
	}

	// 关闭编码器以刷新缓冲区
	err = writer.Close()
	if err != nil {
		fmt.Println("关闭编码器时出错:", err)
		return
	}

	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

// ListAllFilePath  获取所有的文件和目录的全路径，包含子目录下的文件,不包含回车占文件
func ListAllFilePath(dirPath string, dir bool, file bool) []string {
	var pathList []string
	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "$RECYCLE.BIN") {
			if info.IsDir() && !dir {
				return nil
			}
			if !info.IsDir() && !file {
				return nil
			}
			pathList = append(pathList, path)
		}
		return nil
	})
	return pathList
}

// ListAllFileName 获取所有的文件和目录的全路径，包含子目录下的文件,不包含回车占文件
func ListAllFileName(dirPath string, dir bool, file bool) []string {
	var names []string
	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, "$RECYCLE.BIN") {
			if info.IsDir() && !dir {
				return nil
			}
			if !info.IsDir() && !file {
				return nil
			}
			names = append(names, info.Name())
		}
		return nil
	})
	return names
}

// GetRootList 获取所有磁盘的根目录
func GetRootList() []string {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsDrives()
	case "linux":
		return ListFileName("/", "", true, false)
	case "darwin": // macOS 使用的是 darwin
		return ListFileName("/", "", true, false)
	default:
		return nil
	}
}

// GetDBRootFilepath 数据库默认存放路径的根路径
func GetDBRootFilepath() string {
	if runtime.GOOS == "windows" {
		return "D:/myData"
	} else {
		return "/home/myData"
	}
}

func GetWindowsDrives() []string {
	var drives []string
	for letter := 'D'; letter <= 'Z'; letter++ {
		drive := string(letter) + ":"
		if _, err := os.Stat(filepath.Join(drive, ".")); err == nil {
			drives = append(drives, drive)
		}
	}
	return drives
}

// ListFileName 只列出当前文件夹下的文件，不包含子目录中的文件，可以通过后缀限定文件列表
func ListFileName(dirPath string, suffix string, dir bool, file bool) []string {
	// 打开目录
	d, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return nil
	}
	defer Close(d)

	// 读取目录中的文件和子目录
	entries, err := d.ReadDir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	var nameList []string
	// 遍历目录中的文件和子目录
	for _, entry := range entries {
		// 检查是否是文件且后缀匹配
		if entry.IsDir() && !dir {
			continue
		}
		if !entry.IsDir() && !file {
			continue
		}
		if len(suffix) != 0 && !strings.HasSuffix(entry.Name(), suffix) {
			continue
		}
		nameList = append(nameList, entry.Name())
	}
	return nameList
}

// ListFileNameReturnDirAndFile 只列出当前文件夹下的文件，不包含子目录中的文件，可以通过后缀限定文件列表
func ListFileNameReturnDirAndFile(dirPath string) ([]fs.FileInfo, []fs.FileInfo) {
	// 打开目录
	d, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return nil, nil
	}
	defer Close(d)

	// 读取目录中的文件和子目录
	entries, err := d.ReadDir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil, nil
	}

	var dirList []fs.FileInfo
	var nameList []fs.FileInfo
	// 遍历目录中的文件和子目录
	for _, entry := range entries {
		// 检查是否是文件且后缀匹配
		info, _ := entry.Info()
		if entry.IsDir() {
			dirList = append(dirList, info)
		}
		if !entry.IsDir() {
			nameList = append(nameList, info)
		}
	}
	return dirList, nameList
}

// CountFilesInDir 获取目录下的文件数量
func CountFilesInDir(dir string) int {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			count++
		}
	}
	return count
}

func ReadLines(filename string) *[]string {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer Close(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // 读取每一行并添加到切片中
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	return &lines
}

func DiffFiles(file1, file2 string) *[]string {
	return Difference(ReadLines(file1), ReadLines(file2))
}

// CopyFile 复制文件 src-源文件 dst-目标文件
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		log.Printf("文件打开异常:%s,源文件名:%s", err, srcName)
		return
	}
	defer Close(src)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("文件打开异常:%s目标文件名:%s", err, dstName)
		return
	}
	defer Close(dst)
	return io.Copy(dst, src)
}

// CopyDir 递归拷贝整个目录
func CopyDir(src, dst string) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 遍历目录内容
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// 递归拷贝子目录
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// 拷贝文件
			if _, err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CalculateCRC(filePath string) string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 初始化CRC32哈希值
	var crc32Value uint32
	if crc32Value, err = CalculateCRC32(file); err != nil {
		fmt.Printf("Error calculating CRC32: %v\n", err)
		return ""
	}
	return strconv.FormatUint(uint64(crc32Value), 10) // 第二个参数是基数，10表示十进制
}

// CalculateCRC32 calculateCRC32 计算文件的CRC32值
func CalculateCRC32(reader io.Reader) (uint32, error) {
	var hash uint32
	hasher := crc32.NewIEEE() // 使用IEEE版本的CRC32

	// 读取文件内容并更新哈希值
	if _, err := io.Copy(hasher, reader); err != nil {
		return 0, err
	}

	// 获取最终的哈希值
	hash = hasher.Sum32()
	return hash, nil
}

// CleanFileName 清楚文件名字符串的非法字符,以及{"(", ")", "&"}
func CleanFileName(input string) string {
	// 示例字符串，包含一些不允许的文件名字符

	// 定义一个正则表达式，匹配不允许的文件名字符
	// 注意：这里使用了字符类 [] 来包含所有不允许的字符，并使用了转义字符 \\ 来表示字面量的反斜杠
	re := regexp.MustCompile(`[\\/:*?"<>|\s]`)

	// 使用正则表达式替换所有匹配的字符为空字符串
	output := re.ReplaceAllString(input, "")

	invalidChar := []string{"(", ")", "&"}
	for _, str := range invalidChar {
		output = strings.ReplaceAll(output, str, "")
	}
	// 输出处理后的字符串
	return output
}

// ReadFile2Map 读取文件到Map,分隔符
func ReadFile2Map(filePath string, split string) map[string]string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer Close(file)
	// 创建一个scanner来逐行读取文件
	scanner := bufio.NewScanner(file)
	// 初始化map来存储数据
	data := make(map[string]string)
	// 逐行读取文件
	for scanner.Scan() {
		// 获取当前行的文本
		line := scanner.Text()
		// 按冒号分割键和值
		parts := strings.SplitN(line, split, 2)
		if len(parts) == 2 {
			// 去除键和值两边的空白字符
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// 将键值对添加到map中
			data[key] = value
		}
	}
	return data
}

func GetFileMD5(pathName string) string {
	f, err := os.Open(pathName)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}
	defer Close(f)

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}
	has := md5hash.Sum(nil)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// CalculateHash calculateHash 计算文件的MD5哈希值
func CalculateHash(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return ""
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

// CalculateMD5 计算字符串的MD5哈希值
func CalculateMD5(text string) string {
	// 将字符串转换为字节切片
	hash := md5.Sum([]byte(text))
	// 将字节切片转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}

// GetDirSize 获取目录大小，返回格式化字符串
func GetDirSize(dirPath string) string {
	var totalSize int64
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return FormatSize(totalSize)
}

func RenameFile(from string, to string) {
	err := os.Rename(from, to)
	if err != nil {
		fmt.Printf("文件移动出错,from:%s,to:%s\n", from, to)
	}
}

// RemoveBeforeFirstChinese 修改文件名，去除第一个中文字符之前的所有字符
func RemoveBeforeFirstChinese(oldName, basePath string) error {
	index := FindFirstChineseRuneIndex(oldName)
	if index == -1 {
		// 如果没有找到中文字符，则不修改文件名
		fmt.Println("No Chinese character found in the filename:", oldName)
		return nil
	}

	// 截取第一个中文字符及其之后的部分作为新文件名
	newName := oldName[index:]
	// 构建完整的文件路径
	newPath := basePath + "/" + newName

	// 重命名文件或目录
	err := os.Rename(basePath+"/"+oldName, newPath)
	if err != nil {
		return err
	}
	fmt.Printf("Renamed '%s' to '%s'\n", oldName, newName)
	return nil
}

func HasFileByPrefix(v string) bool {
	fnList := ListFileName(ExecPath(), "", false, true)
	for _, fn := range fnList {
		if strings.HasPrefix(fn, v) {
			return true
		}
	}
	return false
}

// FileNameModify 去除吧符合字符串中不能作为文件名的字符
func FileNameModify(input string) string {
	// 示例字符串，包含一些不允许的文件名字符

	// 定义一个正则表达式，匹配不允许的文件名字符
	// 注意：这里使用了字符类 [] 来包含所有不允许的字符，并使用了转义字符 \\ 来表示字面量的反斜杠
	re := regexp.MustCompile(`[\\/:*?"<>|\s]`)

	// 使用正则表达式替换所有匹配的字符为空字符串
	output := re.ReplaceAllString(input, "")

	// 输出处理后的字符串
	return output
}
