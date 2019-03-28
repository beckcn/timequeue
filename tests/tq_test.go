package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"
	
	"github.com/beckcn/timequeue"
)

type HeartBeatData struct {
	uid		uint64
	time	int64
}

func (hb *HeartBeatData) GetKey() string {
	return strconv.FormatUint(hb.uid, 10)
}

func (hb *HeartBeatData) GetValue() int64 {
	return hb.time
}

func checkTimeout(ts *timequeue.TimeQueue) bool {
	if ts.Size() == 0 {
		return false
	}
	now := time.Now().Unix()
	for i := 0; i < 100; i++ {
		get, e := ts.PopTimeout(now)
		if !get {
			break
		}
		fmt.Println("pop element, key: ", e.GetKey(), " value: ", e.GetValue())
	}
	return true
}

func TestTq1(t *testing.T) {
	var (
		tq = timequeue.NewTimeQueue()
		uidStart = uint64(10000)
	)

	insertStart := time.Now().UnixNano() / 1e6
	for i := uint64(0); i < 200000; i++ {
		e := &HeartBeatData{uidStart + i, time.Now().Unix()}
		tq.Push(e)
	}
	t.Log("insert 200000 element spend: ", time.Now().UnixNano() / 1e6 - insertStart, "ms")
}

func TestTq2(t *testing.T) {
	var (
		tq = timequeue.NewTimeQueue()
		uidStart = uint64(10000)
	)

	for i := uint64(0); i < 20; i++ {
		e := &HeartBeatData{uidStart + i, time.Now().Unix()}
		tq.Push(e)
	}

	cb := func(e timequeue.Element) {
		t.Log("key: ", e.GetKey(), " ,value: ", e.GetValue())
	}
	tq.Walk(cb)
}

func TestTq3(t *testing.T) {
	var (
		tq = timequeue.NewTimeQueue()
		uidStart = uint64(10000)
	)

	for i := uint64(0); i < 2000; i++ {
		e := &HeartBeatData{uidStart + i, time.Now().Unix()}
		tq.Push(e)
	}

	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			if !checkTimeout(tq) {
				timer.Stop()
				return
			}
		}
	}
}
