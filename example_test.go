package cache_test
import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"gopkg.in/orivil/cache.v0"
)

var dir = "./testdata"

type data struct {
	Name string
}

func ExampleJsonCache() {
	// 1. new cache struct and new JsonCache
	d := &data{}
	cache, err := cache.NewJsonCache(dir, "data.json")
	if err != nil {
		log.Fatal(err)
	}

	// 2. read not exist file
	// if dir not exist, it will be auto generated
	err = cache.Read(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d.Name == "")

	// 3. cache data to file
	// it will generate a new file
	d.Name = "foobar"
	err = cache.Write(d)
	if err != nil {
		log.Fatal(err)
	}

	// 4. read the file
	d = &data{}
	cache.Read(d)
	fmt.Println(d.Name == "foobar")

	// this test has generated a file, remove test file "./testdata/data.json" after 10 second
	time.Sleep(10 * time.Second)
	os.Remove(filepath.Join(dir, "data.json"))

	// Output:
	// true
	// true
}