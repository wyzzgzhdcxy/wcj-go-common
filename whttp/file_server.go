package whttp

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("文件服务器\n用法: /file?path=文件完整路径\n例: /file?path=E:/images/photo.jpg"))
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
	http.ServeFile(w, r, cleanPath)
}
