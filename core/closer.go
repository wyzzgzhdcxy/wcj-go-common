package core

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

// Close 统一关闭各种资源
func Close(resource interface{}) {
	if closer, ok := resource.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			log.Printf("close err:%s", err)
		}
		return
	}

	// 处理特殊类型
	switch v := resource.(type) {
	case *websocket.Conn:
		if err := v.Close(); err != nil {
			log.Printf("ws.Close() exception:%s", err)
		}
	case *io.PipeReader:
		if err := v.Close(); err != nil {
			fmt.Printf("PipeReader close error:%s", err)
		}
	case *io.PipeWriter:
		if err := v.Close(); err != nil {
			fmt.Printf("PipeWriter close error:%s", err)
		}
	case *io.ReadCloser:
		if err := (*v).Close(); err != nil {
			log.Printf("ReadCloser close error:%s", err)
		}
	default:
		log.Printf("Unsupported resource type for closing:%T", resource)
	}
}
