package myLeveldb

import (
	"fmt"
	"strings"
	"wcj-go-common/core"

	"github.com/syndtr/goleveldb/leveldb"
	levelDbUtil "github.com/syndtr/goleveldb/leveldb/util"
)

func OpenDBFile(filepath string) *leveldb.DB {
	if core.FileExist(filepath) && core.FileExist(filepath+"/CURRENT") && core.FileExist(filepath+"/LOCK") {
		return CreateAndOpenDBFile(filepath)
	} else {
		panic("数据库不存在,路径错误:" + filepath)
	}
}

func CreateAndOpenDBFile(filepath string) *leveldb.DB {
	db, err := leveldb.OpenFile(filepath, nil)
	if err != nil {
		fmt.Printf("filepath:%s,创建leveldb数据库异常!%v", filepath, err.Error())
	}
	return db
}

// PutByte 约定不允许空值存在
func PutByte(db *leveldb.DB, key []byte, value []byte) {
	if len(key) == 0 || len(value) == 0 {
		fmt.Printf("PutByte error, key or value is nil.%s-%s\n", string(key), string(value))
	} else {
		_ = db.Put(key, value, nil)
	}
}

func Put(db *leveldb.DB, key string, value []byte) {
	PutByte(db, []byte(key), value)
}

func PutString(db *leveldb.DB, key string, value string) {
	PutByte(db, []byte(key), []byte(value))
}

func Get(db *leveldb.DB, key string) []byte {
	data, _ := db.Get([]byte(key), nil)
	return data
}

func GetString(db *leveldb.DB, key string) string {
	return string(Get(db, key))
}

func GetMapResult(db *leveldb.DB, key string) *map[string]string {
	var myMap = make(map[string]string)
	myMap[key] = GetString(db, key)
	return &myMap
}

func DeleteKey(db *leveldb.DB, key string) {
	err := db.Delete([]byte(key), nil)
	if err != nil {
		fmt.Printf("删除,KEY:%s异常！%s\n", key, err.Error())
	}
}

func DeleteKeyList(db *leveldb.DB, keyList *[]string) {
	for _, key := range *keyList {
		err := db.Delete([]byte(key), nil)
		if err != nil {
			fmt.Printf("删除,KEY:%s异常！%s\n", key, err.Error())
		}
	}
}

// ListData 查询记录,prefix-前缀 limit 限制条数 keyLen key的长度
func ListData(db *leveldb.DB, prefix string, limit int, keyLen string, value string, valLen string, containKey bool, containValue bool) (int, *[]string, *map[string]string) {
	var myMap = make(map[string]string)
	var iter = db.NewIterator(nil, nil)
	var keys []string
	if core.StrIsNotEmpty(prefix) {
		iter = db.NewIterator(levelDbUtil.BytesPrefix([]byte(prefix)), nil)
	}
	i := 1
	for iter.Next() {
		key := string(iter.Key())
		if len(keyLen) != 0 && !core.AssertStrLen(key, keyLen) {
			continue
		}
		val := string(iter.Value())
		if len(valLen) != 0 && !core.AssertStrLen(val, valLen) {
			continue
		}
		if core.StrIsNotEmpty(value) && !strings.Contains(val, value) {
			continue
		}

		if containValue {
			myMap[key] = val
		}
		if containKey {
			keys = append(keys, key)
		}
		i++
		if limit != 0 && i > limit {
			break
		}
	}
	return i, &keys, &myMap
}

func List(db *leveldb.DB, prefix string, limit int, keyLen string, value string, valLen string) (*[]string, *map[string]string) {
	_, keys, myMap := ListData(db, prefix, limit, keyLen, value, valLen, true, true)
	return keys, myMap
}

func ListKey(db *leveldb.DB, prefix string, limit int, keyLen string, value string, valLen string) *[]string {
	_, keys, _ := ListData(db, prefix, limit, keyLen, value, valLen, true, false)
	return keys
}

func CountKey(db *leveldb.DB, prefix string, keyLen string, value string, valLen string) int {
	count, _, _ := ListData(db, prefix, 0, keyLen, value, valLen, true, false)
	return count
}

type RequestInfo struct {
	TraceId    string
	RequestUrl string
}

func KeyList(db *leveldb.DB) []RequestInfo {
	var requestInfos []RequestInfo
	iter := db.NewIterator(&levelDbUtil.Range{Start: []byte("2022032008324")}, nil)
	var processorStr string
	for iter.Next() {
		key := string(iter.Key())[0:22]
		if !strings.Contains(processorStr, key) {
			url := "[" + key[0:14] + "]" + string(Get(db, key+"_3"))
			requestInfos = append(requestInfos, RequestInfo{TraceId: key, RequestUrl: url})
			processorStr = processorStr + key
		}
	}
	return requestInfos
}

func MultiGet(db *leveldb.DB, keys []string) *map[string]string {
	// 创建一个 map 来存储查询结果
	results := make(map[string]string)
	// 遍历 key 列表，依次查询每个 key
	for _, key := range keys {
		results[key] = GetString(db, key)
	}
	return &results
}

func CloseDB(dbList ...*leveldb.DB) {
	for _, db := range dbList {
		err := db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func CopyDb(fromDb *leveldb.DB, targetDb *leveldb.DB) {
	var iter = fromDb.NewIterator(nil, nil)
	i := 1
	for iter.Next() {
		PutByte(targetDb, iter.Key(), iter.Value())
		i++
	}
	fmt.Printf("数据库复制完成，共拷贝数据记录:%d\n", i)
}
