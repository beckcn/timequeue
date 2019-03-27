package timequeue

import "container/list"

type Element interface {
	GetKey() string
	GetValue() int64
}

type TimeQueue struct {
	dataMap  map[string]*list.Element
	dataList *list.List
}

func NewTimeQueue() *TimeQueue {
	return &TimeQueue{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
	}
}

func (tq *TimeQueue) Exists(e Element) bool {
	_, exists := tq.dataMap[e.GetKey()]
	return exists
}

func (tq *TimeQueue) Push(e Element) bool {
	if tq.Exists(e) {
		tq.dataList.Remove(tq.dataMap[e.GetKey()])
	}
	elem := tq.dataList.PushBack(e)
	tq.dataMap[e.GetKey()] = elem
	return true
}

func (tq *TimeQueue) Size() int {
	return tq.dataList.Len()
}

func (tq *TimeQueue) Walk(cb func(e Element)) {
	for elem := tq.dataList.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value.(Element))
	}
}

func (tq *TimeQueue) PopTimeout(now int64) (bool, Element) {
	front := tq.dataList.Front()
	if front == nil {
		return false, nil
	}
	if e := front.Value.(Element); e.GetValue() <= now {
		tq.dataList.Remove(front)
		delete(tq.dataMap, e.GetKey())
		return true, e
	}
	return false, nil
}