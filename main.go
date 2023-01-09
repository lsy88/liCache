package main

import (
	"fmt"
	"github.com/lsy88/liCache/liCache"
)

func main() {
	memCache := liCache.NewMemCache()
	memCache.SetMaxMemory("200MB")
	
	memCache.Set("zxhang", 10, 10)
	//memCache.Set("li", 20, 10)
	get, _ := memCache.Get("zxhang")
	fmt.Println(get)
	
}
