package cache_test

import (
	"testing"
	"math/rand"
	"strconv"
	"gopkg.in/orivil/cache.v0"
)

var datas = map[int]string{
	1:"1",
	2:"2",
	3:"3",
	4:"4",
	5:"5",
	6:"6",
	7:"7",
	8:"8",
	9:"9",
	10:"10",
}

var cacher = cache.New()

func TestCache(t *testing.T) {
	for id, ins := range datas {
		cacher.Add(id, ins)
	}

	// test GetNext
	order := true // Asc
	results := cacher.GetNext(1, 10, order)
	if len(results) != 10 {
		t.Errorf("length expect:%d, got:%d", 10, len(results))
	}

	pass :=
	results[0].(string) == "1" &&
	results[1].(string) == "2" &&
	results[2].(string) == "3" &&
	// ...
	results[9].(string) == "10"

	if !pass {
		t.Error("asc cacher.GetNext: result unexpected\n")
	}

	// test GetNext
	order = false // Desc
	results = cacher.GetNext(1, 10, order)
	if len(results) != 10 {
		t.Errorf("length expect:%d, got:%d", 10, len(results))
	}

	pass =
	results[0].(string) == "10" &&
	results[1].(string) == "9" &&
	results[2].(string) == "8" &&
	// ...
	results[9].(string) == "1"

	if !pass {
		t.Error("desc cacher.GetNext: result unexpected\n")
	}

	// test GetPrev
	// Asc
	results = cacher.GetPrev(10, 10, true)
	if len(results) != 10 {
		t.Errorf("length expect:%d, got:%d", 10, len(results))
	}

	pass =
	results[0].(string) == "1" &&
	results[1].(string) == "2" &&
	results[2].(string) == "3" &&
	// ...
	results[9].(string) == "10"

	if !pass {
		t.Error("asc cacher.GetPrev: result unexpected\n")
	}

	// test GetPrev
	// Desc
	results = cacher.GetPrev(10, 10, false)
	if len(results) != 10 {
		t.Errorf("length expect:%d, got:%d", 10, len(results))
	}
	pass =
	results[0].(string) == "10" &&
	results[1].(string) == "9" &&
	results[2].(string) == "8" &&
	// ...
	results[9].(string) == "1"

	if !pass {
		t.Error("desc cacher.GetPrev: result unexpected\n")
	}

	// test GetIn
	results = cacher.GetIn([]int{1, 9, 7, 3, 13, 11}, func(id int)interface{} {

		return strconv.Itoa(id)
	})

	if len(results) != 6 {
		t.Errorf("length expect:%d, got:%d", 6, len(results))
	}

	pass =
	results[0].(string) == "1" &&
	results[1].(string) == "9" &&
	results[2].(string) == "7" &&
	results[3].(string) == "3" &&
	results[4].(string) == "13" &&
	results[5].(string) == "11"

	if !pass {
		t.Error("GetIn.GetPrev: result unexpected\n")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 1; i < b.N; i++ {
		cacher.Add(rand.Intn(i), i)
	}
}

func BenchmarkGetNext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		datas := cacher.GetNext(1, 10, true)
		if len(datas) != 10 {
			b.Error("error lenght")
		}
	}
}

func BenchmarkGetPrev(b *testing.B) {
	for i := 0; i < b.N; i++ {
		datas := cacher.GetPrev(100, 10, true)
		if len(datas) != 10 {
			b.Error("error lenght")
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	for i := 1; i < b.N; i++ {
		// if data not cached, it will add to cache
		cacher.Update(rand.Intn(i), i)
	}
}

func BenchmarkDel(b *testing.B) {

	for i := 1; i < b.N; i++ {
		cacher.Del(rand.Intn(i))
	}
}

func BenchmarkGetIn(b *testing.B) {

	var ids = make([]int, 10)

	for i := 1; i < b.N; i++ {
		// init 10 ids
		for j:=0; j < 10; j++ {
			ids[j] = rand.Intn(i)
		}

		// if data not cached, it will be created
		datas := cacher.GetIn(ids, func(id int)interface{} {

			return rand.Intn(i)
		})

		if len(datas) != 10 {
			b.Error("error lenght")
		}

	}
}

//PASS
//BenchmarkAdd-8    	 2000000	       662 ns/op
//BenchmarkGetNext-8	 2000000	       578 ns/op
//BenchmarkGetPrev-8	 3000000	       613 ns/op
//BenchmarkUpdate-8 	10000000	       987 ns/op
//BenchmarkDel-8    	 5000000	       297 ns/op
//BenchmarkGetIn-8  	 1000000	      2779 ns/op
//ok  	gopkg.in/orivil/cache.v0	37.353s
