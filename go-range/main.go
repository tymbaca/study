package main

import "fmt"

type T struct {
	i int
}

type TI interface{}

func NewTI(c <-chan string) TI {
	t := T{}
	c <- "df"
	return t
}

func main() {
	t := NewTI()
	fmt.Printf("%#v", t)
}
