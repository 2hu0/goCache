package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 自定义哈希函数：直接把 key 的数值当作哈希值，方便验证
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// 添加 3 个节点，每个有 3 个虚拟节点：
	// 节点 "2" -> 虚拟节点 hash("02")=2, hash("12")=12, hash("22")=22
	// 节点 "4" -> 虚拟节点 hash("04")=4, hash("14")=14, hash("24")=24
	// 节点 "6" -> 虚拟节点 hash("06")=6, hash("16")=16, hash("26")=26
	// 哈希环上排序后：2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",  // hash=2, 命中虚拟节点 2 -> 节点 "2"
		"11": "2",  // hash=11, 顺时针第一个 >=11 的是 12 -> 节点 "2"
		"23": "4",  // hash=23, 顺时针第一个 >=23 的是 24 -> 节点 "4"
		"27": "2",  // hash=27, 比 26 大，回绕到环起点 2 -> 节点 "2"
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 添加节点 "8"，虚拟节点：8, 18, 28
	// 哈希环变为：2, 4, 6, 8, 12, 14, 16, 18, 22, 24, 26, 28
	hash.Add("8")

	// 27 原本映射到 "2"，现在应该映射到 "8"（28 是第一个 >=27 的虚拟节点）
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}