package circuit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketLinkList(t *testing.T) {

	m := &BucketLinklist{}
	assert.Equal(t, true, m.IsEmpty())
	assert.Equal(t, float32(0.0), m.CalcFailRate(10))

	m.RemoveHead()
	assert.Nil(t, m.head)

	m.UpsertLast(&SecondMetric{
		Ts: time.Now().Unix(),
	})
	assert.Equal(t, float32(0.0), m.CalcFailRate(1))

	m.RemoveHead()
	assert.Nil(t, m.head)

	now1 := time.Now().Unix()
	m.UpsertLast(&SecondMetric{
		Ts:   now1,
		Req:  10,
		Fail: 6,
	})

	assert.Equal(t, 10, m.tail.val.Req)
	assert.Equal(t, 6, m.tail.val.Fail)
	assert.Equal(t, now1, m.tail.val.Ts)

	time.Sleep(time.Second)
	now2 := time.Now().Unix()
	m.UpsertLast(&SecondMetric{
		Ts:   now2,
		Req:  6,
		Fail: 1,
	})

	assert.Equal(t, 6, m.tail.val.Req)
	assert.Equal(t, 1, m.tail.val.Fail)
	assert.Equal(t, now2, m.tail.val.Ts)

	m.UpsertLast(&SecondMetric{
		Ts:   now2,
		Req:  10,
		Fail: 5,
	})
	assert.Equal(t, 16, m.tail.val.Req)
	assert.Equal(t, 6, m.tail.val.Fail)
	assert.Equal(t, now2, m.tail.val.Ts)

	m.UpsertLast(&SecondMetric{
		Ts:   now2 - 1,
		Req:  10,
		Fail: 5,
	})
	assert.Equal(t, 16, m.tail.val.Req)
	assert.Equal(t, 6, m.tail.val.Fail)
	assert.Equal(t, now2, m.tail.val.Ts)

	time.Sleep(time.Second)

	now3 := time.Now().Unix()
	m.UpsertLast(&SecondMetric{
		Ts:   now3,
		Req:  10,
		Fail: 5,
	})

	assert.Equal(t, 10, m.tail.val.Req)
	assert.Equal(t, 5, m.tail.val.Fail)
	assert.Equal(t, now3, m.tail.val.Ts)

	// timeIdx reqCount failCount interval
	// 1       10        6         1
	// 2       6         1         1
	// 2       10        5         0
	// 3       10        5         1
	// ---------------------------
	// rate last 1s = 5 / 10
	// rate last 2s = 11 / 26
	// rate last 3s = 17 / 36
	// rate last ns = 17 / 36
	assert.Equal(t, float32(0.5), m.CalcFailRate(1))
	assert.Equal(t, float32(11.0/26.0), m.CalcFailRate(2))
	assert.Equal(t, float32(17.0/36.0), m.CalcFailRate(3))
	assert.Equal(t, float32(17.0/36.0), m.CalcFailRate(10))

	assert.Equal(t, 3, m.Size())

	m.RemoveHead()
	assert.Equal(t, 2, m.Size())

	// b := m.Buckets()
	// assert.Equal(t, now2, b.next.val.Ts)
	// assert.Equal(t, now3, b.next.val.Ts)

	assert.Equal(t, float32(0.5), m.CalcFailRate(1))
	assert.Equal(t, float32(11.0/26.0), m.CalcFailRate(4))
	assert.Equal(t, float32(11.0/26.0), m.CalcFailRate(20))

}

func TestBuckets(t *testing.T) {
	m := &BucketLinklist{}
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 1, len(m.Buckets()))
}

func TestCalcFailRate(t *testing.T) {
	m := &BucketLinklist{}
	now1 := time.Now().Unix()
	m.UpsertLast(&SecondMetric{
		Ts:   now1,
		Req:  10,
		Fail: 6,
	})
	m.RemoveHead()

	m.UpsertLast(&SecondMetric{
		Ts:   now1,
		Req:  10,
		Fail: 5,
	})
	m.UpsertLast(&SecondMetric{
		Ts:   now1,
		Req:  10,
		Fail: 5,
	})
	assert.Equal(t, float32(0.5), m.CalcFailRate(1))

	time.Sleep(time.Second)
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	time.Sleep(time.Second)
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 3, m.Size())
	assert.Equal(t, float32(6.0/10.0), m.CalcFailRate(1))
	assert.Equal(t, float32(12.0/20.0), m.CalcFailRate(2))
	assert.Equal(t, float32(22.0/40.0), m.CalcFailRate(3))

	m.RemoveHead()
	assert.Equal(t, 2, m.Size())
	assert.Equal(t, float32(6.0/10.0), m.CalcFailRate(1))
	assert.Equal(t, float32(12.0/20.0), m.CalcFailRate(2))
	assert.Equal(t, float32(12.0/20.0), m.CalcFailRate(3))

	m.RemoveHead()
	assert.Equal(t, 1, m.Size())
	assert.Equal(t, float32(6.0/10.0), m.CalcFailRate(1))
	assert.Equal(t, float32(6.0/10.0), m.CalcFailRate(2))
	assert.Equal(t, float32(6.0/10.0), m.CalcFailRate(3))
}

func TestGetLastQPS(t *testing.T) {
	m := &BucketLinklist{}
	now := time.Now().Unix()
	assert.Equal(t, 0, m.GetLastQPS(now))
	m.UpsertLast(&SecondMetric{
		Ts:   now,
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 10, m.GetLastQPS(now))
	m.UpsertLast(&SecondMetric{
		Ts:   now,
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 20, m.GetLastQPS(now))
	assert.Equal(t, 0, m.GetLastQPS(now+1000))
}

func TestMonitor1(t *testing.T) {
	m := NewBucketLinklist(1)
	m.monitor(1)
	time.Sleep(time.Second)
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  0,
		Fail: 6,
	})
	assert.Equal(t, float32(1.0), m.CalcFailRate(1))
}

func TestMonitor(t *testing.T) {
	m := NewBucketLinklist(1)
	m.monitor(1)
	time.Sleep(time.Second)

	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 1, m.Size())

	time.Sleep(time.Second)
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 2, m.Size())

	time.Sleep(time.Second)
	m.UpsertLast(&SecondMetric{
		Ts:   time.Now().Unix(),
		Req:  10,
		Fail: 6,
	})
	assert.Equal(t, 2, m.Size())
	time.Sleep(time.Second)
	assert.Equal(t, 1, m.Size())
}
