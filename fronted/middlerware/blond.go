package middlerware

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
)

var base = 1 << 32

// achieve blond filter
// 1. calculate the hash of key
// 2. preload the players data
var (
	Client *redis.Client
	Prefix = "product:"
)

type IBloomFilter interface {
	Contains(productID int64) bool
	Add(productID int64) error
	Print()
}

type BloomFilter struct {
	n    int
	p    float64
	size int     //位图的大小
	k    int     //k个哈希函数
	bits *bitmap //Redis的bitmap
	//根据k个值获得k个种子数组
	SEEDS []int
	//根据k个种子，创建k个哈希函数并添加到数组中
	funcs []*SimpleHash
}

// 初始化函数
func NewBloomFilter(n int, p float64) IBloomFilter {
	size := getSizeOfBloomFilter(n, p)
	k := getNumberOfHashFunc(size, n)
	SEEDS := getKDistinctPrimes(k)
	funcs := make([]*SimpleHash, 0) 
	// 根据 SEEDS 中的每个种子创建 SimpleHash 实例并添加到 funcs 切片中

	for _, seed := range SEEDS {
		hash := NewSimpleHash(size, seed)
		funcs = append(funcs, hash)
	}
	bitmaps := NewBitMap(size)
	return &BloomFilter{n, p, size, k, bitmaps, SEEDS, funcs}
}
func (bf *BloomFilter) Add(productID int64) error {
	for _, f := range bf.funcs {
		hashValue := int(f.hash(productID))
		bf.bits.set(hashValue)

	}
	return nil
}

func (bf *BloomFilter) Contains(productID int64) bool {
	for _, f := range bf.funcs {
		hashValue := int(f.hash(productID))
		if !bf.bits.has(hashValue) {
			// 如果任何一个哈希函数计算的位置上为 0，则说明该元素一定不存在于过滤器中
			return false
		}
	}
	// 如果所有哈希函数计算的位置上都为 1，则说明该元素可能存在于过滤器中
	return true
}

func (bf *BloomFilter) Print() {
	fmt.Println(bf.n, bf.k, bf.p, bf.size)
}
func getKDistinctPrimes(k int) []int {
	primes := make([]int, 0)
	num := 2

	for len(primes) < k {
		if isPrime(num) && isWellDistributed(num, primes) {
			primes = append(primes, num)
		}
		num++
	}

	return primes
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	if num == 2 || num == 3 {
		return true
	}
	if num%2 == 0 || num%3 == 0 {
		return false
	}

	for i := 5; i*i <= num; i += 6 {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
	}

	return true
}

func isWellDistributed(num int, primes []int) bool {
	for _, prime := range primes {
		if math.Abs(float64(prime-num)) < 20 { // 设置一个距离阈值，确保质数之间距离较远
			return false
		}
	}
	return true
}

type SimpleHash struct {
	cap  int
	seed int
}

func NewSimpleHash(cap, seed int) *SimpleHash {
	return &SimpleHash{
		cap:  cap,
		seed: seed,
	}
}

func (s *SimpleHash) hash(value int64) int64 {
	h := int64(0)
	highBits := value >> 32
	lowBits := value & 0xFFFFFFFF
	if value != 0 {
		h = int64(math.Abs(float64(s.seed * (s.cap - 1) & (int(highBits) ^ int(lowBits)))))
	}
	return h
}

type bitmap struct {
	keys []byte
	len  int
}

func NewBitMap(len int) *bitmap {
	return &bitmap{keys: make([]byte, len/8+1), len: len}
}

func (b *bitmap) has(v int) bool {
	k := v / 8
	kv := byte(v % 8)
	if k > len(b.keys) { //todo not exist
		return false
	}
	if b.keys[k]&(1<<kv) != 0 {
		return true
	}
	return false
}

func (b *bitmap) set(v int) {
	k := v / 8        // 计算元素 v 对应的字节索引
	kv := byte(v % 8) // 计算元素 v 在字节中的位偏移
	for b.len <= k {  // 如果位图长度小于等于 k，则说明当前位图中还没有包含索引为 k 的字节，需要扩容
		b.keys = append(b.keys, 0) // 在位图的末尾添加一个新的字节，初始值为 0
		b.len++
	}
	b.keys[k] = b.keys[k] | (1 << kv) // 将索引为 k 的字节中的第 kv 位设置为 1，表示元素 v 存在于位图中
}

func (b *bitmap) length() int {
	return b.len
}

// 根据传入的数据量 n 和 误差率 p 返回 size
func getSizeOfBloomFilter(n int, p float64) int {
	size := int(math.Ceil(-(float64(n) * math.Log(p)) / math.Pow(math.Log(2), 2)))
	return size
}

// 根据传入的 size 和 数据量 n 返回所需的哈希函数值
func getNumberOfHashFunc(size, n int) int {
	k := int(math.Ceil((float64(size) / float64(n)) * math.Log(2)))
	return k
}
