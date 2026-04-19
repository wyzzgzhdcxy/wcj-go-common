package core

import (
	"fmt"
	"hash/fnv"
	"wcj-go-common/model"
)

func MapExistKey[T any](dstFilesMap *map[string]T, key string) bool {
	_, ok := (*dstFilesMap)[key]
	return ok
}

// 定义一个函数来计算键的哈希值，并根据哈希值返回对应的索引（0 到 numPartitions-1）
func hashToPartition(key string, numPartitions int) int {
	h := fnv.New64a()
	_, err := h.Write([]byte(key))
	if err != nil {
		fmt.Println(err.Error())
	}
	hashValue := h.Sum64()
	return int(hashValue % uint64(numPartitions))
}

func SplitMap[T comparable](inputMap *map[string]T, targetMapSize int) *[]map[string]T {
	mapLen := len(*inputMap)
	numPartitions := mapLen / targetMapSize
	if mapLen%targetMapSize != 0 {
		numPartitions++
	}
	smallMaps := make([]map[string]T, numPartitions)
	if mapLen <= targetMapSize {
		smallMaps[0] = *inputMap
	} else {
		// 定义要切分成的小 map 的数量
		// 创建一个切片来存储小 map
		for i := range smallMaps {
			smallMaps[i] = make(map[string]T)
		}
		// 遍历大 map，并根据哈希值将键值对分配到小 map 中
		for key, value := range *inputMap {
			partitionIndex := hashToPartition(key, numPartitions)
			smallMaps[partitionIndex][key] = value
		}
		// 打印结果
	}
	return &smallMaps
}

func Map2EntityList(keys *[]string, myMap *map[string]string) *[]model.MapEntity {
	var mapEntityList []model.MapEntity
	for _, k := range *keys {
		mapEntityList = append(mapEntityList, model.MapEntity{
			Key:   k,
			Value: (*myMap)[k],
		})
	}
	return &mapEntityList
}

// MapMerge 通用合并函数（Go 1.18+ 泛型）
// 类似java hashmap的merge函数
func MapMerge[K comparable, V any](m *map[K]V, key K, newVal V, fn func(V, V) V) {
	if oldVal, ok := (*m)[key]; ok {
		(*m)[key] = fn(oldVal, newVal)
	} else {
		(*m)[key] = newVal
	}
}

// MapMergeInt 通用合并函数（Go 1.18+ 泛型）
// 类似java hashmap的merge函数
func MapMergeInt(m *map[string]int, key string, newVal int) {
	if oldVal, ok := (*m)[key]; ok {
		(*m)[key] = oldVal + newVal
	} else {
		(*m)[key] = newVal
	}
}
