package main

import (
	"fmt"

	"github.com/tymbaca/study/myhashmap/hashmap"
)

func main() {
	hm := hashmap.New()

	fmt.Println(hm)
	hm.Set("a", "name")
	fmt.Println(hm)
	hm.Set("b", "name")
	hm.Set("a", "another")
	fmt.Println(hm)
}
