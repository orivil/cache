package cache_test
import (
	"log"
	"fmt"
	"gopkg.in/orivil/cache.v0"
)

var dir = "./testdata"

type data struct {
	Name string
}

func ExampleRead() {
	d := &data{Name: "golang"}

	// 1. new JsonCache
	cache, err := cache.NewJsonCache(dir, "data.json")
	if err != nil {
		log.Fatal(err)
	}

	// 2. read exist or not exist file
	// if dir not exist, it will be auto generated
	err = cache.Read(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d.Name == "foobar")

	// Output:
	// true
}

func ExampleWrite() {
	d := &data{Name: "foobar"}

	// 1. new JsonCache
	cache, err := cache.NewJsonCache(dir, "data.json")
	if err != nil {
		log.Fatal(err)
	}

	// 2. cache data to file
	// it will cover the exist file or generate a new file
	err = cache.Write(d)
	if err != nil {
		log.Fatal(err)
	}
}