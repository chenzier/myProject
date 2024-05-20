package datamodels

import (
	"fmt"
	"sync"
	"time"
)

// 号段
type LeafSegment struct {
	Cursor uint64 // 当前发放位置
	Max    uint64 // 最大值
	Min    uint64 // 开始值即最小值
	InitOk bool   // 是否初始化成功
}

// 以初始化DB中的MAXID为1举例子，也就是号段从1开始，步长为1000,最大值就是1000，范围是1～1000。下一段范围就是1001～2000
func NewLeafSegment(leaf *Leaf) *LeafSegment {
	return &LeafSegment{
		Cursor: leaf.MaxID - uint64(leaf.Step+1), // 最小值的前一个值
		Max:    leaf.MaxID - 1,                   // DB默认存的是1 所以这里要减1
		Min:    leaf.MaxID - uint64(leaf.Step),   // 开始的最小值
		InitOk: true,
	}
}

type LeafAlloc struct {
	Key        string                 // 也就是`biz_tag`用来区分业务
	Step       int32                  // 记录步长
	CurrentPos int32                  // 当前使用的 segment buffer光标; 总共两个buffer缓存区，循环使用
	Buffer     []*LeafSegment         // 双buffer 一个作为预缓存作用
	UpdateTime time.Time              // 记录更新时间 方便长时间不用进行清理，防止占用内存
	mutex      sync.Mutex             // 互斥锁
	IsPreload  bool                   // 是否正在预加载
	Waiting    map[string][]chan byte // 挂起等待
}

type Leaf struct {
	ID          uint64 `json:"id" form:"id"`                   // 主键id
	BizTag      string `json:"biz_tag" form:"biz_tag"`         // 区分业务
	MaxID       uint64 `json:"max_id" form:"max_id"`           // 该biz_tag目前所被分配的ID号段的最大值
	Step        int32  `json:"step" form:"step"`               // 每次分配ID号段长度
	Description string `json:"description" form:"description"` // 描述
	UpdateTime  uint64 `json:"update_time" form:"update_time"` // 更新时间
}

func NewLeafAlloc(leaf *Leaf) *LeafAlloc {
	return &LeafAlloc{
		Key:        leaf.BizTag,
		Step:       leaf.Step,
		CurrentPos: 0, // 初始化使用第一块buffer缓存
		Buffer:     make([]*LeafSegment, 0),
		UpdateTime: time.Now(),
		Waiting:    make(map[string][]chan byte), //初始化
		IsPreload:  false,
	}
}

func (l *LeafAlloc) Lock() {
	l.mutex.Lock()
}

func (l *LeafAlloc) Unlock() {
	l.mutex.Unlock()
}

func (l *LeafAlloc) HasSeq() bool {
	if l.Buffer[l.CurrentPos].InitOk && l.Buffer[l.CurrentPos].Cursor < l.Buffer[l.CurrentPos].Max {
		return true
	}
	return false
}

func (l *LeafAlloc) HasID(id uint64) bool {
	return id != 0
}

func (l *LeafAlloc) Wakeup() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for _, waitChan := range l.Waiting[l.Key] {
		close(waitChan)
	}
	l.Waiting[l.Key] = l.Waiting[l.Key][:0]
}

const ExpiredTime = time.Minute * 15 //清理超过15min没更新的缓存

// 全局分配器
// key: biz_tag value: SegmentBuffer
type LeafSeq struct {
	cache sync.Map
}

func NewLeafSeq() *LeafSeq {
	seq := &LeafSeq{}
	go seq.clear()
	return seq
}

// 获取
func (l *LeafSeq) Get(bizTag string) *LeafAlloc {
	if seq, ok := l.cache.Load(bizTag); ok {
		return seq.(*LeafAlloc)
	}
	return nil
}

// 添加
func (l *LeafSeq) Add(seq *LeafAlloc) string {
	l.cache.Store(seq.Key, seq)
	return seq.Key
}

// 更新
func (l *LeafSeq) Update(key string, bean *LeafAlloc) {
	if element, ok := l.cache.Load(key); ok {
		alloc := element.(*LeafAlloc)
		alloc.Buffer = bean.Buffer
		alloc.UpdateTime = bean.UpdateTime
	}
}

// 清理超过15min没用过的内存
func (l *LeafSeq) clear() {
	for {
		now := time.Now()
		// 15分钟后
		mm, _ := time.ParseDuration("15m")
		next := now.Add(mm)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		fmt.Println("start clear goroutine")
		l.cache.Range(func(key, value interface{}) bool {
			alloc := value.(*LeafAlloc)
			if next.Sub(alloc.UpdateTime) > ExpiredTime {
				fmt.Printf("clear biz_tag: %s cache", key)
				l.cache.Delete(key)
			}
			return true
		})
	}
}
