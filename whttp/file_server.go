package whttp

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// OnActivate 是单例激活回调，由调用方注册。
// 当 HTTP 服务收到激活请求时会被调用。
var OnActivate func()

var allowedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
	".webp": true, ".svg": true, ".ico": true, ".tiff": true, ".tif": true,
	".heic": true, ".heif": true, ".psd": true, ".raw": true, ".dng": true,
	".cr2": true, ".nef": true, ".arw": true, ".jp2": true, ".j2k": true,
	".icns": true, ".wbmp": true, ".jxr": true, ".hdp": true, ".wdp": true,
	".mp3": true, ".wav": true, ".flac": true, ".ogg": true, ".aac": true,
	".m4a": true, ".wma": true, ".ape": true, ".opus": true, ".alac": true,
	".ac3": true, ".dts": true, ".aiff": true, ".aif": true, ".mid": true,
	".midi": true, ".oga": true, ".spx": true, ".amr": true, ".mmf": true,
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".webm": true,
	".flv": true, ".wmv": true, ".m4v": true, ".3gp": true, ".mpeg": true,
	".mpg": true, ".mpv": true, ".ts": true, ".m2ts": true, ".mts": true,
	".vob": true, ".ogv": true, ".rm": true, ".rmvb": true, ".qt": true,
	".swf": true, ".f4v": true, ".smil": true,
}

func Start(port string) (addr string, err error) {
	addr = "127.0.0.1" + port

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("端口 %s 被占用: %v", port, err)
	}

	http.HandleFunc("/file", fileHandler)
	http.HandleFunc("/audio", audioHandler)
	http.HandleFunc("/activate", activateHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("文件服务器\n用法:\n  /file?path=文件完整路径   (通用文件，支持图片/音频/视频)\n  /audio?path=音频路径     (音频流，支持拖动/进度)\n  /activate               (单例激活信号)\n例: /file?path=E:/images/photo.jpg"))
	})

	go func() {
		http.Serve(ln, nil)
	}()

	return addr, nil
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	var filePath string

	if r.URL.RawQuery != "" {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		filePath = params.Get("path")
	}

	if filePath == "" {
		http.Error(w, "请提供 path 参数\n例: /file?path=E:/images/photo.jpg", http.StatusBadRequest)
		return
	}

	filePath = strings.ReplaceAll(filePath, "\\", "/")
	cleanPath := filepath.Clean(filePath)

	if strings.Contains(cleanPath, "..") {
		http.Error(w, "非法路径", http.StatusForbidden)
		return
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "文件不存在: "+cleanPath, http.StatusNotFound)
		} else {
			http.Error(w, "无法访问: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		http.Error(w, "暂不支持目录访问", http.StatusForbidden)
		return
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))
	if !allowedExts[ext] {
		http.Error(w, "不支持的文件类型，仅允许图片、音频、视频文件", http.StatusForbidden)
		return
	}

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Range")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Range")
	http.ServeFile(w, r, cleanPath)
}

var audioExts = map[string]bool{
	".mp3": true, ".wav": true, ".flac": true, ".ogg": true, ".aac": true,
	".m4a": true, ".wma": true, ".ape": true, ".opus": true, ".alac": true,
	".aiff": true, ".aif": true, ".mid": true, ".midi": true, ".oga": true,
	".spx": true, ".amr": true, ".mmf": true,
}

func activateHandler(w http.ResponseWriter, r *http.Request) {
	if OnActivate != nil {
		OnActivate()
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func audioContentType(ext string) string {
	switch ext {
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".aac":
		return "audio/aac"
	case ".mp3":
		return "audio/mpeg"
	}
	return "audio/mpeg"
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "请提供 path 参数\n例: /audio?path=E:/music/song.mp3", http.StatusBadRequest)
		return
	}

	filePath = strings.ReplaceAll(filePath, "\\", "/")
	cleanPath := filepath.Clean(filePath)

	if strings.Contains(cleanPath, "..") {
		http.Error(w, "非法路径", http.StatusForbidden)
		return
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "文件不存在: "+cleanPath, http.StatusNotFound)
		} else {
			http.Error(w, "无法访问: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if info.IsDir() {
		http.Error(w, "暂不支持目录访问", http.StatusForbidden)
		return
	}

	ext := strings.ToLower(filepath.Ext(cleanPath))
	if !audioExts[ext] {
		http.Error(w, "不支持的音频格式", http.StatusForbidden)
		return
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		http.Error(w, "无法打开文件: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	contentType := audioContentType(ext)

	// ETag 和缓存头
	etag := fmt.Sprintf(`"%x"`, fileStat.ModTime().UnixNano())
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "private, max-age=3600")

	// 检查 If-None-Match，返回 304
	if cachedETag := r.Header.Get("If-None-Match"); cachedETag != "" && cachedETag == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Range, ETag")

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		io.Copy(w, file)
		return
	}

	// 解析 Range 头，格式: bytes=start-end
	rangePart := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangePart, "-")
	if len(parts) != 2 {
		http.Error(w, "无效的 Range 头", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, _ := strconv.ParseInt(parts[0], 10, 64)
	end := fileSize - 1
	if parts[1] != "" {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	}

	if start > end || start >= fileSize {
		http.Error(w, "Range 不满足", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	contentLength := end - start + 1
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.WriteHeader(http.StatusPartialContent)

	file.Seek(start, 0)
	io.CopyN(w, file, contentLength)
}
