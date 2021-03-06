package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//Hash maps bytes to uint32
type Hash func(data []byte) uint32

//Map contains all hashed keys
type Map struct {
	hash Hash
	replicas int//一个真实节点对应几个虚拟节点
	keys []int//sorted
	hashMap map[int]string
}

//New Creates a Map instance
func New(replicas int,fn Hash)*Map{
	m:= &Map{
		replicas:replicas,
		hash:fn,
		hashMap:make(map[int]string),
	}
	if m.hash == nil{
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//Add some keys to the hash
func(m *Map) Add(keys ...string){
	for _,key := range keys{
		for i:=0;i<m.replicas;i++{
			hash := int(m.hash([]byte(strconv.Itoa(i)+key)))
			m.keys = append(m.keys,hash)
			m.hashMap[hash] = key//虚拟节点对应的真实节点
		}
	}
	sort.Ints(m.keys)
}

//Get gets the closest item in the hash to the provided key
func (m *Map)Get(key string) string{
	if len(m.keys) == 0{
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//Binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i]>=hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}