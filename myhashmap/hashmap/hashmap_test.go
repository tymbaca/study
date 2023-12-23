package hashmap

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_Set(t *testing.T) {
	m := NewHashMap()
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("%d", rand.Intn(60))
		val := fmt.Sprintf("%d", rand.Intn(60))
		m.Set(key, val)
	}
	fmt.Println(m)
}
