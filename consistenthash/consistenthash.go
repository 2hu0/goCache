package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 哈希函数类型，输入字节切片，输出 uint32 哈希值
type Hash func(data []byte) uint32

// Map 一致性哈希环，管理所有节点在环上的分布
type Map struct {
	hash     Hash           // 哈希函数，默认 crc32.ChecksumIEEE
	replicas int            // 每个真实节点对应的虚拟节点数量
	keys     []int          // 哈希环上所有虚拟节点的哈希值，已排序
	hashMap  map[int]string // 虚拟节点哈希值 -> 真实节点名称的映射
}

// New 创建一致性哈希环
// replicas: 每个真实节点对应的虚拟节点数，越大分布越均匀
// fn: 自定义哈希函数，传 nil 则默认使用 crc32.ChecksumIEEE
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 向哈希环添加节点
// 每个真实节点会生成 replicas 个虚拟节点（"{i}{key}" 的哈希值），
// 均匀散布在环上，解决节点过少时的数据倾斜问题
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点的哈希：hash(str(i) + key)，i 不同则哈希值不同
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key // 哈希值映射回真实节点
		}
	}
	sort.Ints(m.keys) // 排序后方便二分查找
}

// Get 根据 key 查找它应该落在哪个节点上
// 在哈希环上顺时针找到第一个 >= hash(key) 的虚拟节点，
// 该虚拟节点对应的真实节点就是 key 所属节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// 二分查找第一个 >= hash 的虚拟节点位置
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 如果 hash 比所有虚拟节点都大（idx == len(m.keys)），
	// 则回绕到环的起点（idx % len == 0），形成环形结构
	return m.hashMap[m.keys[idx%len(m.keys)]]
}