package myBadger

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"github.com/wyzzgzhdcxy/wcj-go-common/core"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/options"
)

func CompactDb(db *badger.DB) {
	// 手动触发Compaction
	// 注意：这里的RunValueLogGC并不是直接触发Compaction，而是触发Value Log的GC（垃圾回收）
	// 它会删除已标记为删除的值，并可能触发底层的Compaction过程。
	// 如果你想要更细粒度的控制，你可能需要直接使用badger的内部API，但这通常不推荐。
	err := db.RunValueLogGC(0.7) // 0.7是GC的阈值，表示当可用空间低于70%时触发GC
	if err != nil {
		log.Fatalf("无法运行Value Log GC: %v", err)
	}

	// 注意：上面的RunValueLogGC并不是直接等价于Compaction，但它会间接影响Compaction过程。
	// Badger的Compaction是自动管理的，并且是由内部逻辑控制的。
	// 如果你确实需要更直接地控制Compaction（例如，为了测试或特定的性能调优），
	// 你可能需要深入研究Badger的源代码，并可能需要对代码进行修改。

	// 在生产环境中，通常不建议这样做，因为Badger的自动Compaction机制已经经过优化，
	// 并且能够适应大多数用例。

	fmt.Println("手动触发了Value Log GC（可能间接影响Compaction过程）")
}

func OpenDBFile(filepath string) *badger.DB {
	if core.FileExist(filepath) && core.FileExist(filepath+"/MANIFEST") && core.FileExist(filepath+"/DISCARD") {
		db, err := badger.Open(getOpts(filepath))
		if err != nil {
			log.Fatalf("无法打开 Badger 数据库: %v", err)
		}
		return db
	} else {
		panic("数据库不存在,路径错误:" + filepath)
	}
}

func NewOpenDBFile(filepath string) (*badger.DB, error) {
	if core.FileExist(filepath) && core.FileExist(filepath+"/MANIFEST") && core.FileExist(filepath+"/DISCARD") {
		return badger.Open(getOpts(filepath))
	} else {
		return nil, errors.New("数据库不存在,路径错误:" + filepath)
	}
}

func CreateAndOpenDBFile(filepath string) *badger.DB {
	db, err := badger.Open(getOpts(filepath))
	if err != nil {
		log.Fatalf("无法打开 Badger 数据库: %v", err)
	}
	return db
}

func CreateDBFile(filepath string) {
	db := CreateAndOpenDBFile(filepath)
	defer db.Close()
}

// PutByte 约定不允许空值存在
func PutByte(db *badger.DB, key []byte, value []byte) {
	if len(key) == 0 || len(value) == 0 {
		fmt.Printf("PutByte error, key or value is nil.%s-%s\n", string(key), string(value))
	} else {
		err := db.Update(func(txn *badger.Txn) error {
			err := txn.Set(key, value)
			if err != nil {
				log.Fatalf("写入数据失败: %v", err)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("写入数据失败: %v", err)
		}
	}
}

func Put(db *badger.DB, key string, value []byte) {
	PutByte(db, []byte(key), value)
}

func PutMap(db *badger.DB, myMap *map[string]string) {
	mapArr := core.SplitMap(myMap, 20000)
	for _, smallMap := range *mapArr {
		err := db.Update(func(txn *badger.Txn) error {
			for key, val := range smallMap {
				err := txn.Set([]byte(key), []byte(val))
				if err != nil {
					log.Fatalf("写入数据失败: %v", err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("写入数据失败: %v", err)
		}
	}
}

func PutString(db *badger.DB, key string, value string) {
	PutByte(db, []byte(key), []byte(value))
}

func Get(db *badger.DB, key string) []byte {
	var val []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
	}
	return val
}

func GetString(db *badger.DB, key string) string {
	return string(Get(db, key))
}

func KeyExists(db *badger.DB, key string) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		return err
	})
	return err == nil
}

func GetMapResult(db *badger.DB, key string) *map[string]string {
	var myMap = make(map[string]string)
	myMap[key] = GetString(db, key)
	return &myMap
}

func DeleteKey(db *badger.DB, key string) {
	// 执行删除操作
	err := db.Update(func(txn *badger.Txn) error {
		// 尝试删除键
		err := txn.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("无法删除键: %w", err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("删除键时出错: %v", err)
	}
}

func DeleteKeyList(db *badger.DB, keyList *[]string) {
	for _, key := range *keyList {
		// 执行删除操作
		err := db.Update(func(txn *badger.Txn) error {
			// 尝试删除键
			err := txn.Delete([]byte(key))
			if err != nil {
				return fmt.Errorf("无法删除键: %w", err)
			}
			return nil
		})

		if err != nil {
			log.Fatalf("删除键时出错: %v", err)
		}
	}
}
func List(db *badger.DB, prefix string, limit int, keyLen string, value string, valLen string) (*[]string, *map[string]string) {
	_, keys, myMap := ListData(db, prefix, limit, keyLen, value, valLen, true, true)
	return keys, myMap
}

// ListData List 查询记录,prefix-前缀 limit 限制条数 keyLen key的长度
func ListData(db *badger.DB, prefix string, limit int, keyLen string, value string, valLen string, containKey bool, containValue bool) (int, *[]string, *map[string]string) {
	var myMap = make(map[string]string)
	var keys []string
	i := 0

	// 遍历数据库
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 如果不需要值，可以设置为false以减少内存使用
		iter := txn.NewIterator(opts)
		defer iter.Close()

		if len(prefix) != 0 {
			iter.Seek([]byte(prefix))
		} else {
			iter.Rewind()
		}

		for ; iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			// 检查键是否以前缀开头
			if len(prefix) != 0 && !strings.HasPrefix(key, prefix) {
				break // 如果不是，则停止迭代
			}

			valByte, err := item.ValueCopy(nil)
			val := string(valByte)
			if err != nil {
				return fmt.Errorf("无法复制值: %w", err)
			}

			if len(keyLen) != 0 && !core.AssertStrLen(key, keyLen) {
				continue
			}
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
			if limit != 0 && i >= limit {
				break
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("遍历数据库时出错: %v", err)
	}

	return i, &keys, &myMap
}

func getOpts(dstPath string) badger.Options {
	//1024*1024=1048576
	var MB int64 = 1048576
	opts := badger.DefaultOptions(dstPath).WithLoggingLevel(badger.ERROR).WithBaseLevelSize(128 * MB).WithBaseTableSize(64 * MB).WithMemTableSize(256 * MB)
	opts.Compression = options.ZSTD
	return opts
}

func ListKey(db *badger.DB, prefix string, limit int, keyLen string, value string, valLen string) *[]string {
	_, keys, _ := ListData(db, prefix, limit, keyLen, value, valLen, true, false)
	return keys
}

func CountKey(db *badger.DB, prefix string, keyLen string, value string, valLen string) uint64 {
	count, _, _ := ListData(db, prefix, 0, keyLen, value, valLen, true, false)
	return uint64(uint(count))
}

func CountKeyFast(db *badger.DB) (uint64, error) {
	count := uint64(0)
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 不读取值以加速
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
		}
		return nil
	})
	return count, err
}

type SumTotalCountInfo struct {
	Count     uint64
	Timestamp int64
}

func isOlderThan10Minutes(info SumTotalCountInfo) bool {
	currentTime := time.Now().Unix() // 当前时间的 Unix 时间戳（秒）
	diff := currentTime - info.Timestamp
	return diff > 10*60 // 10分钟 = 600秒
}

// CountKeyWithCache 数据库缓存命名规则
// meta_sum_total_count     有效时间5分钟,格式:时间戳|count
func CountKeyWithCache(db *badger.DB) (uint64, error) {
	sumTotalCountInfo := Get(db, "meta_sum_total_count")
	if len(sumTotalCountInfo) != 0 {
		totalInfo := SumTotalCountInfo{}
		core.JsonToObject(&sumTotalCountInfo, &totalInfo)
		if !isOlderThan10Minutes(totalInfo) {
			return totalInfo.Count, nil
		}
	}
	count, _ := CountKeyFast(db)
	totalInfo := SumTotalCountInfo{
		Count:     count,
		Timestamp: time.Now().Unix(),
	}
	PutString(db, "meta_sum_total_count", core.ToString(totalInfo))
	return totalInfo.Count, nil
}

func MultiGet(db *badger.DB, keys *[]string) *map[string]string {

	// 存储键和对应值的map
	values := make(map[string]string)

	// 执行读操作
	err := db.View(func(txn *badger.Txn) error {
		for _, key := range *keys {
			item, err := txn.Get([]byte(key))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					// 键不存在，可以记录日志或执行其他操作
					continue
				}
				return fmt.Errorf("无法获取键的值: %w", err)
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("无法复制值: %w", err)
			}
			// 将键和值存储到map中
			values[key] = string(val)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return &values
}

func MultiGetWithPrefix(db *badger.DB, keys *[]string, prefix string) *map[string]string {

	// 存储键和对应值的map
	values := make(map[string]string)

	// 执行读操作
	err := db.View(func(txn *badger.Txn) error {
		for _, key := range *keys {
			item, err := txn.Get([]byte(prefix + key))
			if err != nil {
				if errors.Is(err, badger.ErrKeyNotFound) {
					// 键不存在，可以记录日志或执行其他操作
					continue
				}
				return fmt.Errorf("无法获取键的值: %w", err)
			}
			val, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("无法复制值: %w", err)
			}
			// 将键和值存储到map中
			values[key] = string(val)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return &values
}

func CloseDB(dbList ...*badger.DB) {
	for _, db := range dbList {
		err := db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func CopyDbWithCallback(srcPath string, dstPath string, callback func([]byte)) {
	CopyDbWithPrefix(srcPath, dstPath, nil, callback)
}

func CopyDbWithPrefix(srcPath string, dstPath string, prefix []byte, callback func([]byte)) {
	//如果不要添加前缀，切记不要调用这个方法，这个方法性能会下降
	srcDB := OpenDBFile(srcPath)
	dstDB := CreateAndOpenDBFile(dstPath)
	defer CloseDB(srcDB, dstDB)
	CopyDbWithPrefixParmsDb(srcDB, dstDB, prefix, callback)
}

func CopyDbWithPrefixParmsDb(srcDB *badger.DB, dstDB *badger.DB, prefix []byte, callback func([]byte)) {
	//如果不要添加前缀，切记不要调用这个方法，这个方法性能会下降
	length := len(prefix)
	i := 0
	srcDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 1000 // 控制预取量
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				return dstDB.Update(func(txn *badger.Txn) error {
					if length > 0 {
						return txn.Set(append(prefix, item.KeyCopy(nil)...), val)
					} else {
						return txn.Set(item.KeyCopy(nil), val)
					}
				})
			})
			if err != nil {
				return err
			}
			i++
			if i%10000 == 0 {
				callback([]byte(fmt.Sprintf("复制进度:%d\n", i)))
			}
		}
		return nil
	})
	callback([]byte(fmt.Sprintf("数据库复制完成，共拷贝数据记录:%d\n", i)))
}

func GetAll(db *badger.DB) *map[string]string {
	var myMap = make(map[string]string)
	// 遍历数据库
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 如果不需要值，可以设置为false以减少内存使用
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			valByte, _ := item.ValueCopy(nil)
			myMap[key] = string(valByte)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("遍历数据库时出错: %v", err)
	}
	return &myMap
}

func GetAllByte(db *badger.DB) *map[string][]byte {
	var myMap = make(map[string][]byte)
	// 遍历数据库
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // 如果不需要值，可以设置为false以减少内存使用
		iter := txn.NewIterator(opts)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			valByte, _ := item.ValueCopy(nil)
			myMap[key] = valByte
		}
		return nil
	})
	if err != nil {
		log.Fatalf("遍历数据库时出错: %v", err)
	}
	return &myMap
}
