package main

import (
	"fmt"
	"time"
)

type My struct {
	S []byte
}

func main() {
	// h := make([]byte, 40*1024*1024)
	ptrs := make([]*My, 0)
	c := 1
	i := 0
	for {
		s := make([]byte, 1024*1024)
		m := &My{S: s}
		ptrs = append(ptrs, m)
		fmt.Printf("Added struct in ptrs, size: %d\n", len(ptrs))

		if i%10 == 0 {
			ptrs = make([]*My, 0)
			fmt.Printf("Cleared ptrs (%d times)\n", c)
			c++
		}

		time.Sleep(100 * time.Millisecond)
		i++
	}
}
