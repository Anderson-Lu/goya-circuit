package circuit

import (
	"sync"
	"sync/atomic"
	"time"
)

type BucketLinkNode struct {
	pre  *BucketLinkNode
	next *BucketLinkNode
	val  *SecondMetric
}

type BucketLinklist struct {
	sync.Mutex
	head     *BucketLinkNode
	tail     *BucketLinkNode
	qpsMutex sync.RWMutex
	state    uint32
	size     int
}

func NewBucketLinklist(windowSize int) *BucketLinklist {
	m := &BucketLinklist{}
	atomic.StoreUint32(&m.state, 0)
	m.monitor(windowSize)
	return m
}

func (s *BucketLinklist) UpsertLast(metric *SecondMetric) {

	s.Lock()
	defer s.Unlock()

	s.qpsMutex.Lock()
	defer s.qpsMutex.Unlock()

	// list is empty and inset to s.head
	if s.IsEmpty() {
		n := &BucketLinkNode{
			val: metric,
		}
		s.head = n
		s.tail = n
		s.size = 1
		return
	}

	n := &BucketLinkNode{
		pre: s.tail,
		val: metric,
	}

	// invalid time
	if s.tail.val.Ts > n.val.Ts {
		return
	}

	// equals the last node and update it
	if s.tail.val.Ts == n.val.Ts {
		s.tail.val.Fail += n.val.Fail
		s.tail.val.Req += n.val.Req
		return
	}

	n.pre = s.tail
	s.tail.next = n
	s.tail = n
	s.size++
}

func (s *BucketLinklist) RemoveHead() {

	if s.IsEmpty() {
		return
	}

	if s.head.next == nil {
		s.Clear()
		s.size = 0
		return
	}

	s.head = s.head.next
	s.head.pre = nil
	s.size--
}

func (s *BucketLinklist) Size() int {
	return s.size
}

func (s *BucketLinklist) Clear() {
	s.head = nil
	s.tail = nil
}

func (s *BucketLinklist) IsEmpty() bool {
	return s.head == nil
}

func (s *BucketLinklist) Buckets() []SecondMetric {
	var r = make([]SecondMetric, 0)
	var h = s.head
	for h != nil {
		r = append(r, *h.val)
		h = h.next
	}
	return r
}

func (s *BucketLinklist) GetLastQPS(ts int64) int {

	s.qpsMutex.RLock()
	defer s.qpsMutex.RUnlock()

	if s.IsEmpty() {
		return 0
	}

	if s.tail.val.Ts != ts {
		return 0
	}

	return s.tail.val.Req
}

func (s *BucketLinklist) CalcFailRate(windowSize int) float32 {

	if s.IsEmpty() {
		return 0.0
	}

	var failCount = 0
	var reqCount = 0
	var cursor = s.tail
	var now = time.Now().Unix()

	for cursor != nil {
		if now-cursor.val.Ts < int64(windowSize) {
			failCount += cursor.val.Fail
			reqCount += cursor.val.Req
			cursor = cursor.pre
			continue
		}
		break
	}

	if reqCount == 0 && failCount == 0 {
		return 0.0
	}

	// This is because the request has been initiated before,
	// but now it ends, and a timeout has occurred.
	// and there are no new requests in the last window
	// In this scenario, the default value is all failures
	if reqCount == 0 && failCount > 0 {
		return 1.0
	}

	return float32(failCount) / float32(reqCount)
}

func (s *BucketLinklist) monitor(windowSize int) {

	if !atomic.CompareAndSwapUint32(&s.state, 0, 1) {
		return
	}

	go s.tick(windowSize)
}

func (s *BucketLinklist) tick(windowSize int) {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		now := time.Now().Unix()
		if s.IsEmpty() {
			continue
		}
		if now-s.head.val.Ts > int64(windowSize) {
			s.RemoveHead()
			continue
		}
	}
}
