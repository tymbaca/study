package main

import (
	"fmt"
	"time"
)

type T struct {
	l []int
}

type TI interface{}

func main() {
	now := time.Date(2024, time.January, 32, 0, 0, 0, 0, time.UTC)
	fmt.Println(now.String())
	next := now.AddDate(0, 1, 0)
	fmt.Println(next.String())
}
