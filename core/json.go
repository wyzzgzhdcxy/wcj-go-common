package core

import (
	"encoding/json"
	"log"
)

func JsonToObject(content *[]byte, v interface{}) {
	if len(*content) == 0 {
		println("content is empty")
	}
	err := json.Unmarshal(*content, v)
	if err != nil {
		log.Printf("JsonToObject json.Unmarshal err:%v", err)
	}
}

func JsonToObjectWithMsg(content *[]byte, v interface{}, msg string) {
	if len(*content) == 0 {
		println("content is empty")
	}
	err := json.Unmarshal(*content, v)
	if err != nil {
		log.Printf("JsonToObject json.Unmarshal err:%v", err)
		log.Printf("错误信息:%s", msg)
	}
}

func ToJsonString(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Printf("ToJsonString json.Marshal err:%s", err)
		return ""
	} else {
		return string(jsonData)
	}
}

func ToJsonByte(v interface{}) []byte {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Printf("ToJsonString json.Marshal err:%s", err)
		return nil
	} else {
		return jsonData
	}
}

func ToByteArray(v interface{}) []byte {
	jsonData, err := json.Marshal(v)
	if err != nil {
		log.Printf("ToByteArray json.Marshal err:%s", err)
		return nil
	} else {
		return jsonData
	}
}
