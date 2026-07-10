package gocache

// PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter。

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// Get 是从group中查找缓存的值,PeerGetter 就对应于 HTTP 客户端。
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
