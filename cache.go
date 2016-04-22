package cache

import (
	"sync"
)

type trie struct {
	data interface{}
	prev *trie
	next *trie
}

type Cache struct {
	intIndex map[int]*trie
	sync.RWMutex
	maxID    int
}

func New() *Cache {
	return &Cache{
		intIndex:make(map[int]*trie, 1000),
	}
}

func (this *Cache) Del(id int) {
	this.Lock()
	defer this.Unlock()
	if t, ok := this.intIndex[id]; ok {
		if t.prev != nil {
			t.prev.next = t.next
		}
		if t.next != nil {
			t.next.prev = t.prev
		}
		delete(this.intIndex, id)
	}
}

func (this *Cache) Update(id int, ins interface{}) {
	this.RLock()
	t, ok := this.intIndex[id]
	this.RUnlock()
	if ok {
		// ignore data race
		t.data = ins
	} else {
		this.Add(id, ins)
	}
}

// if has id not exist in cache, it will call 'dataProvider', and cache the result
func (this *Cache) GetIn(ids []int, dataProvider func(id int)interface{}) (datas []interface{}) {
	this.RLock()
	datas = make([]interface{}, len(ids))
	for i, id := range ids {
		if t, ok := this.intIndex[id]; ok {
			datas[i] = t.data
		} else {
			this.RUnlock()
			// data not exist in cache, get it form 'dataProvider'
			data := dataProvider(id)
			// add to cache
			datas[i] = data
			this.Add(id, data)
			this.RLock()
		}
	}
	this.RUnlock()
	return
}

func (this *Cache) GetNext(startID, num int, orderAsc bool) (datas []interface{}) {
	datas = make([]interface{}, num)

	this.RLock()
	defer this.RUnlock()
	maxID := this.maxID
	if startID > maxID {
		return nil
	} else if startID < 1 {
		startID = 1
	}

	for startID <= maxID {
		if t, ok := this.intIndex[startID]; ok {
			// asc
			if orderAsc {
				index := 0
				for {
					if num > index {
						datas[index] = t.data
						if t.next == nil {
							return datas[0:index + 1]
						}
						t = t.next
						index++
					} else {
						return
					}
				}
			// desc
			} else {
				index := num - 1
				for {
					if index >= 0 {
						datas[index] = t.data
						if t.next == nil {
							return datas[index:]
						}
						t = t.next
						index--
					} else {
						return
					}
				}
			}
		} else {
			startID++
		}
	}

	return nil
}

func (this *Cache) GetPrev(startID, num int, orderDesc bool) (datas []interface{}) {
	datas = make([]interface{}, num)

	this.RLock()
	defer this.RUnlock()
	if startID > this.maxID {
		startID = this.maxID
	} else if startID < 1 {
		// get the last data
		startID = this.maxID
	}
	for startID > 0 {
		if t, ok := this.intIndex[startID]; ok {
			// asc
			if !orderDesc {
				index := 0
				for {
					if num > index {
						datas[index] = t.data
						if t.prev == nil {
							return datas[0:index + 1]
						}
						t = t.prev
						index++
					} else {
						return
					}
				}
			// desc
			} else {
				index := num - 1
				for {
					if index >= 0 {
						datas[index] = t.data
						if t.prev == nil {
							return datas[index:]
						}
						t = t.prev
						index--
					} else {
						return
					}
				}
			}
		} else {
			startID--
		}
	}

	return nil
}

func (this *Cache) Add(id int, ins interface{}) {
	var prevID int
	var nextID int
	var toTop bool
	var toButtom bool

	this.Lock()
	defer this.Unlock()
	if _, ok := this.intIndex[id]; !ok {
		if this.maxID < id {
			this.maxID = id
		}
		t := &trie{data: ins}
		this.intIndex[id] = t
		index := 1
		for !toTop || !toButtom {
			prevID = id - index
			nextID = id + index
			// add the trie to prev
			if prevID > 0 {
				if prevT, ok := this.intIndex[prevID]; ok {
					if prevT.next != nil {
						prevT.next.prev = t
						t.next = prevT.next
					}
					prevT.next = t
					t.prev = prevT
					break
				}
			} else {
				toTop = true
			}
			// add the trie to next
			if nextID <= this.maxID {
				if nextT, ok := this.intIndex[nextID]; ok {
					if nextT.prev != nil {
						nextT.prev.next = t
						t.prev = nextT.prev
					}
					nextT.prev = t
					t.next = nextT
					break
				}
			} else {
				toButtom = true
			}
			// add both sides
			index++
		}
	}
}
